// 來源：bodhi-guide 首頁的解說員（guide-hero-assets 分支 guide-hero/avatar.js）。
// 官網版多了 dispose()，離開 AI 組長頁時停止繪製並釋放 WebGL 資源。
// 覺行小組線上組長 3D 角色：站在及膝的湖水中冥想，雙手始終捧著一杯熱茶。
// 被提問時睜眼、鏡頭推近到臉部；回答時跟著語音對嘴，身體只輕微擺動。回答結束後回到冥想、鏡頭拉遠。
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { VRMLoaderPlugin, VRMUtils } from '@pixiv/three-vrm';
import { buildScene, WATER_Y, STAND_Z, SHORE_Z, REFLECT_LAYER } from './scene.js';
import { buildDuskScene } from './scene-dusk.js';
import { buildMeadowScene } from './scene-meadow.js';
import { buildPinesScene } from './scene-pines.js';

const VOWELS = ['aa', 'ih', 'ou', 'ee', 'oh'];
const CHARS_PER_SEC = 7;          // 沒有聲音時（靜音或裝置不支援）依字幕逐字對嘴的速度
const SYLLABLES_PER_SEC = 5.5;    // 有聲音時嘴型開合的速度，約等於華語每秒唸的字數
const FACE_SPAN = 0.45;
const CLOSE_FRAC = Number(new URLSearchParams(location.search).get('closeFrac') || 0.5);   // 草原場景推近時，鏡頭往 Sunny 移動的比例           // 近景時畫面寬度至少要容下的範圍（公尺），窄的手機畫面鏡頭會退後一點
const PAUSE = /[\s，。、；：！？,.;:!?「」『』（）()…—\n]/;

// 手臂姿勢（各骨骼在「身體座標」中指向的方向；面向 +Z，角色的左手邊是 +X）：
// 雙手前臂往前、往中線收，掌心相對，在胸前捧著茶杯
const ARMS = {
  leftUpperArm:  [0.16, -0.93, 0.33],
  rightUpperArm: [-0.16, -0.93, 0.33],
  leftLowerArm:  [-0.28, 0.25, 0.93],
  rightLowerArm: [0.28, 0.25, 0.93],
};
const ARM_TWIST = { leftLowerArm: -Math.PI / 2, rightLowerArm: Math.PI / 2 };   // 沿前臂轉 90°，掌心相對、拇指朝上
const ARM_BONES = Object.keys(ARMS);
// 捧杯：四指順著杯身彎曲（右手正、左手負），兩隻拇指在杯子後上方彎成弧、指尖相碰
const GRIP = {
  IndexProximal: 0.6, IndexIntermediate: 0.6, IndexDistal: 0.4,
  MiddleProximal: 0.65, MiddleIntermediate: 0.6, MiddleDistal: 0.4,
  RingProximal: 0.7, RingIntermediate: 0.6, RingDistal: 0.4,
  LittleProximal: 0.75, LittleIntermediate: 0.6, LittleDistal: 0.4,
};
const THUMB = { ThumbMetacarpal: [0.45, 0.1], ThumbProximal: [0.35, 0.3], ThumbDistal: [0.2, 0.2] };   // [往前彎, 往掌心收]
const CUP_OFFSET = new THREE.Vector3(0, 0, 0.025);   // 杯子中心相對兩手腕中點（身體座標）
const restDir = name => new THREE.Vector3(...(name.startsWith('left') ? [1, 0, 0] : [-1, 0, 0]));

const smoothstep = (a, b, x) => { const t = Math.min(1, Math.max(0, (x - a) / (b - a))); return t * t * (3 - 2 * t); };
const angleLerp = (a, b, k) => { let d = ((b - a + Math.PI) % (Math.PI * 2)) - Math.PI; if (d < -Math.PI) d += Math.PI * 2; return a + d * k; };

export async function createAvatar(canvas, url, { onProgress, onIdle, time = 'day' } = {}) {
  const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  renderer.outputColorSpace = THREE.SRGBColorSpace;

  const scene = new THREE.Scene();
  const camera = new THREE.PerspectiveCamera(35, 1, 0.05, 1000);
  const env = ({ dusk: buildDuskScene, meadow: buildMeadowScene, pines: buildPinesScene }[time] ?? buildScene)(scene, camera, renderer);   // 白天藍天湖景、黃昏湖景、晨霧草原或松林雪山

  const loader = new GLTFLoader();
  loader.register(parser => new VRMLoaderPlugin(parser));
  const gltf = await loader.loadAsync(url, e => e.total && onProgress?.(e.loaded / e.total));
  const vrm = gltf.userData.vrm;
  VRMUtils.removeUnnecessaryVertices(gltf.scene);
  VRMUtils.rotateVRM0(vrm);
  vrm.scene.traverse(o => { o.frustumCulled = false; o.layers.enable(REFLECT_LAYER); });   // 人物也映在水面上
  scene.add(vrm.scene);
  const S = env.avatarScale ?? 1;   // 草原場景依全景圖的比例把人物縮小
  vrm.scene.scale.setScalar(S);

  const bone = name => vrm.humanoid.getNormalizedBoneNode(name);
  const B = Object.fromEntries(['hips', 'spine', 'chest', 'neck', 'head',
    'leftUpperLeg', 'rightUpperLeg', 'leftLowerLeg', 'rightLowerLeg', ...ARM_BONES].map(n => [n, bone(n)]));
  for (const [side, sign] of [['right', 1], ['left', -1]]) {
    for (const [name, curl] of Object.entries(GRIP)) bone(side + name)?.rotation.set(0, 0, sign * curl);
    for (const [name, [x, y]] of Object.entries(THUMB)) bone(side + name)?.rotation.set(x, sign * y, 0);
  }
  const hands = ['leftHand', 'rightHand'].map(n => vrm.humanoid.getRawBoneNode(n));

  // 依「身體座標」中的目標方向擺手臂：扣掉整個角色的轉向，求出骨骼的區域旋轉
  const qBodyInv = new THREE.Quaternion(), qParent = new THREE.Quaternion(), qTarget = new THREE.Quaternion(), qTwist = new THREE.Quaternion();
  const aimBody = (name, dirBody, twist = 0) => {
    const b = B[name]; if (!b) return;
    b.parent.getWorldQuaternion(qParent).premultiply(qBodyInv);
    qTarget.setFromUnitVectors(restDir(name), dirBody);
    if (twist) qTarget.multiply(qTwist.setFromAxisAngle(restDir(name), twist));
    b.quaternion.copy(qParent.invert().multiply(qTarget));
    b.updateMatrixWorld(true);
  };

  const cup = makeTeaCup();
  cup.group.traverse(o => o.layers.enable(REFLECT_LAYER));
  cup.group.scale.setScalar(S);
  scene.add(cup.group);
  const cupPos = new THREE.Vector3();

  // 站進水裡：讓水面剛好在膝蓋（小腿骨的起點）高度；走上岸時沿著坡度出水
  vrm.scene.position.set(0, 0, STAND_Z);
  vrm.scene.updateMatrixWorld(true);
  const kneeY = B.leftLowerLeg.getWorldPosition(new THREE.Vector3()).y;
  const depthAt = z => env.onGround ? kneeY * (env.wadeDepth ?? 0) : kneeY * (1 - smoothstep(SHORE_Z - 0.3, SHORE_Z + 0.2, z));   // 草原場景站在河裡，水深依場景設定
  vrm.scene.position.y = WATER_Y - depthAt(STAND_Z);
  vrm.update(0);
  vrm.scene.updateMatrixWorld(true);

  // 蘆葦最高只到肩膀（上臂骨的起點）
  env.addVegetation(B.leftUpperArm.getWorldPosition(new THREE.Vector3()).y, WATER_Y - kneeY - 0.01);   // 湖底就在腳底
  const headOffset = B.head.getWorldPosition(new THREE.Vector3()).y - vrm.scene.position.y;   // 頭部相對腳底的高度

  // 視線目標（睜眼時看鏡頭）
  const gaze = new THREE.Object3D();
  scene.add(gaze);
  if (vrm.lookAt) vrm.lookAt.target = gaze;

  // 風吹頭髮：調整彈簧骨的重力方向
  const joints = [...(vrm.springBoneManager?.joints ?? [])];
  const windDir = new THREE.Vector3(0.8, -0.3, 0.5).normalize();

  const reduceMotion = matchMedia('(prefers-reduced-motion: reduce)').matches;
  const em = vrm.expressionManager;
  const set = (name, v) => em?.getExpression(name) && em.setValue(name, v);
  const ease = (cur, target, k, dt) => cur + (target - cur) * Math.min(1, k * dt);

  let state = 'idle', pendingIdle = false, voiced = false;   // voiced：語音正在唸
  const queue = [];
  const mouth = { vowel: 'aa', amount: 0, target: 0 };
  const face = { eyesClosed: 1, happy: 0, relaxed: 0 };
  let nextBlink = 3, blinkT = -1, shot = 0, talk = 0;   // shot: 0 = 遠景，1 = 近景；talk: 回答時的輕微擺動

  const body = { x: env.standAt?.x ?? 0, z: env.standAt?.z ?? STAND_Z, yaw: 0 };   // 草原場景指定 Sunny 站的位置

  const resize = () => {
    const w = canvas.clientWidth, h = canvas.clientHeight;
    if (!w || !h) return;
    renderer.setSize(w, h, false);
    camera.aspect = w / h;
    // 鏡頭平移（不傾斜）：把畫面往上挪，角色落在中上方，下方留給字幕與輸入框
    camera.setViewOffset(w, h, 0, Math.round(h * 0.14), w, h);
    camera.updateProjectionMatrix();
    env.resize(w * renderer.getPixelRatio(), h * renderer.getPixelRatio());
  };
  const resizeObs = new ResizeObserver(resize);
  resizeObs.observe(canvas);
  resize();

  // 捲到頁面下方、看不到場景時暫停繪製
  let onScreen = true;
  // 草原場景：按住滑鼠左鍵（或手指橫向拖曳）旋轉鏡頭環看四周，放開幾秒後慢慢轉回原本的構圖
  const look = { yaw: 0, pitch: 0, dragging: false, lastX: 0, lastY: 0, idle: 0 };
  const onDown = e => {
    if (e.button !== 0) return;
    look.dragging = true; look.lastX = e.clientX; look.lastY = e.clientY;
    canvas.setPointerCapture(e.pointerId); canvas.style.cursor = 'grabbing';
  };
  const onMove = e => {
    if (!look.dragging) return;
    const k = Math.PI / canvas.clientWidth * zoom.cur;   // 拖過整個畫面寬度約轉 180°（拉近時轉得少一些）
    look.yaw += (e.clientX - look.lastX) * k;
    look.pitch = THREE.MathUtils.clamp(look.pitch + (e.clientY - look.lastY) * k, -0.6, 0.6);
    look.lastX = e.clientX; look.lastY = e.clientY; look.idle = 0;
  };
  const onUp = e => {
    if (!look.dragging) return;
    look.dragging = false; canvas.style.cursor = 'grab';
    if (canvas.hasPointerCapture(e.pointerId)) canvas.releasePointerCapture(e.pointerId);
  };
  if (env.wide) {
    canvas.style.cursor = 'grab';
    canvas.style.touchAction = 'pan-y';   // 直向滑動仍然捲動頁面，橫向拖曳才轉鏡頭
    canvas.addEventListener('pointerdown', onDown);
    canvas.addEventListener('pointermove', onMove);
    canvas.addEventListener('pointerup', onUp);
    canvas.addEventListener('pointercancel', onUp);
  }
  const lookDir = new THREE.Vector3(), lookRight = new THREE.Vector3();

  // 滑鼠滾輪縮放：往上滾拉近（朝游標所指的地方）、往下滾拉遠回原本的構圖；
  // 已經是原本大小時再往下滾，就照常捲動頁面
  const ZOOM_MIN = 0.3;   // 最多拉近到原本視角的 0.3 倍（約 3 倍望遠）
  const zoom = { target: 1, cur: 1 };
  const onWheel = e => {
    const dy = e.deltaMode === 1 ? e.deltaY * 16 : e.deltaY;   // 以行為單位的滾輪換成像素
    if (dy > 0 && zoom.target >= 1) return;                    // 沒有拉近時往下滾：讓頁面捲動
    e.preventDefault();
    const before = zoom.target;
    zoom.target = THREE.MathUtils.clamp(zoom.target * Math.exp(dy * 0.0015), ZOOM_MIN, 1);
    // 朝游標方向拉近：鏡頭跟著轉一點，讓游標下的景物大致留在原處
    if (env.wide && zoom.target < before) {
      const r = canvas.getBoundingClientRect();
      const nx = ((e.clientX - r.left) / r.width) * 2 - 1, ny = 1 - ((e.clientY - r.top) / r.height + 0.14) * 2;   // 0.14：畫面往上挪的量（見 resize）
      const tanV = Math.tan(THREE.MathUtils.degToRad(camera.fov / 2)), shift = 1 - zoom.target / before;
      look.yaw -= Math.atan(nx * tanV * camera.aspect) * shift;
      look.pitch = THREE.MathUtils.clamp(look.pitch + Math.atan(ny * tanV) * shift, -0.6, 0.6);
    }
    look.idle = 0;
  };
  canvas.addEventListener('wheel', onWheel, { passive: false });
  const baseFov = camera.fov;

  const visibleObs = new IntersectionObserver(([e]) => { onScreen = e.isIntersecting; });
  visibleObs.observe(canvas);

  const camPos = new THREE.Vector3(), camLook = new THREE.Vector3(), wideLook = new THREE.Vector3(), closeLook = new THREE.Vector3();
  const widePos = new THREE.Vector3(), closePos = new THREE.Vector3(), dirBody = new THREE.Vector3(), hipsW = new THREE.Vector3();
  // 除錯用：?cam=x,y,z,lookX,lookY,lookZ 固定鏡頭
  const debugCam = new URLSearchParams(location.search).get('cam')?.split(',').map(Number);
  if (debugCam || new URLSearchParams(location.search).has('debug')) window.__bodhi = { vrm, THREE, body, camera, env, ARMS, ARM_TWIST, CUP_OFFSET, setState: s => { state = s; } };
  if (env.wide) camera.position.copy(env.wide.pos);
  else camera.position.set(0.5, vrm.scene.position.y + headOffset + 0.25, 6.4);   // 從遠景開始，避免第一幀從原點飛進來
  const clock = new THREE.Clock();
  let t = 0, charClock = 0;

  renderer.setAnimationLoop(() => {
    const dt = Math.min(clock.getDelta(), 0.1);
    if (!onScreen) return;
    t += dt;
    env.uTime.value = t;
    const m = reduceMotion ? 0.3 : 1;
    env.update(t, dt, m);   // 游魚、天鵝、燕子

    // ---- 站在原地，慢慢轉身面向鏡頭 ----
    const faceYaw = Math.atan2(camera.position.x - body.x, camera.position.z - body.z);
    body.yaw = angleLerp(body.yaw, faceYaw, Math.min(1, 2.5 * dt));
    const depth = depthAt(body.z);
    vrm.scene.position.set(body.x, WATER_Y - depth, body.z);
    vrm.scene.rotation.y = body.yaw;
    vrm.scene.updateMatrixWorld(true);
    qBodyInv.copy(vrm.scene.quaternion).invert();

    // ---- 身體、頭、手臂：回答時身體只輕微搖動 ----
    talk = ease(talk, state === 'talking' ? 1 : 0, 2, dt);
    const breath = Math.sin(t * (state === 'idle' ? 1.1 : 1.6));
    if (B.chest) B.chest.rotation.x = -0.03 + breath * 0.02 * m;
    if (B.spine) B.spine.rotation.set(0.04 + Math.sin(t * 0.7) * 0.012 * talk * m, Math.sin(t * 0.5) * 0.03 * talk * m, Math.sin(t * 0.9) * 0.015 * talk * m);
    vrm.scene.updateMatrixWorld(true);
    for (const name of ARM_BONES) aimBody(name, dirBody.set(...ARMS[name]).normalize(), ARM_TWIST[name]);
    let hx = state === 'idle' ? 0.12 : 0.02, hy = 0, hz = 0;
    if (state === 'thinking') { hz = 0.08; hx = 0.05; }
    if (state === 'talking') { hx += Math.sin(t * 2.3) * 0.02 * m; hy = Math.sin(t * 1.3) * 0.03 * m; }
    B.head.rotation.x = ease(B.head.rotation.x, hx, 3, dt);
    B.head.rotation.y = ease(B.head.rotation.y, hy, 3, dt);
    B.head.rotation.z = ease(B.head.rotation.z, hz, 3, dt);

    // ---- 身旁的蘆葦往兩旁分開 ----
    vrm.scene.updateMatrixWorld(true);
    env.setPushers([{ position: B.hips.getWorldPosition(hipsW), radius: 0.5 }]);
    env.setRipple(body.x, body.z, depth / kneeY);

    // ---- 表情：冥想閉眼、思考半閉、回答睜眼帶笑 ----
    const target = state === 'talking' ? { eyesClosed: 0, happy: 0.35, relaxed: 0 }
      : state === 'thinking' ? { eyesClosed: 0.45, happy: 0, relaxed: 0.3 }
      : { eyesClosed: 0.92, happy: 0, relaxed: 0.25 };
    for (const k of Object.keys(face)) face[k] = ease(face[k], target[k], 2.5, dt);
    set('happy', face.happy); set('relaxed', face.relaxed);
    let blink = face.eyesClosed;
    if (face.eyesClosed < 0.3) {          // 睜眼時才眨眼
      if (blinkT < 0 && t > nextBlink) blinkT = 0;
      if (blinkT >= 0) {
        blinkT += dt;
        blink = Math.max(blink, blinkT < 0.07 ? blinkT / 0.07 : Math.max(0, 1 - (blinkT - 0.07) / 0.1));
        if (blinkT > 0.17) { blinkT = -1; nextBlink = t + 2.5 + Math.random() * 3.5; }
      }
    }
    set('blink', blink);

    // ---- 對嘴：有聲音時跟著語音開合；沒有聲音時依字幕逐字對嘴 ----
    charClock += dt;
    if (voiced) {
      if (charClock >= 1 / SYLLABLES_PER_SEC) {
        charClock = 0;
        if (Math.random() < 0.15) mouth.target = 0.05;
        else { mouth.vowel = VOWELS[Math.floor(Math.random() * VOWELS.length)]; mouth.target = 0.35 + Math.random() * 0.5; }
      }
    } else if (charClock >= 1 / CHARS_PER_SEC) {
      charClock = 0;
      const ch = queue.shift();
      if (ch === undefined || PAUSE.test(ch)) mouth.target = 0;
      else { mouth.vowel = VOWELS[ch.codePointAt(0) % VOWELS.length]; mouth.target = 0.45 + (ch.codePointAt(0) % 7) / 14; }
    }
    mouth.amount = ease(mouth.amount, mouth.target, 18, dt);
    for (const v of VOWELS) set(v, v === mouth.vowel ? mouth.amount : ease(em?.getValue(v) ?? 0, 0, 18, dt));
    if (state === 'talking' && !voiced && queue.length === 0 && mouth.amount < 0.02 && pendingIdle) { state = 'idle'; pendingIdle = false; onIdle?.(); }

    // ---- 風吹頭髮 ----
    const gust = 0.06 + 0.05 * (Math.sin(t * 1.3) * 0.6 + Math.sin(t * 2.6) * 0.25 + 0.35) * m;
    for (const j of joints) { j.settings.gravityDir.copy(windDir); j.settings.gravityPower = gust; }

    // ---- 鏡頭：冥想時遠景；被提問、回答時推近到臉部，並緩慢漂移 ----
    const headY = vrm.scene.position.y + headOffset;
    if (env.wide) {   // 草原場景：鏡頭在全景圖拍攝點；直式手機畫面窄，鏡頭往 Sunny 那邊轉，她才不會擠在邊上
      widePos.copy(env.wide.pos);
      const k = THREE.MathUtils.clamp((1.3 - camera.aspect) / 0.8, 0, 1);
      const toSunny = Math.atan2(body.x - widePos.x, widePos.z - body.z);
      wideLook.copy(env.wide.look).sub(widePos).applyAxisAngle(THREE.Object3D.DEFAULT_UP, -toSunny * 0.85 * k).add(widePos);
    }
    else {
      widePos.set(0.5 + body.x * 0.35, headY + 0.25, 6.4);
      wideLook.set(body.x, headY - 0.1, body.z - 1.5);
    }
    const near = THREE.MathUtils.clamp(FACE_SPAN / (2 * Math.tan(THREE.MathUtils.degToRad(camera.fov / 2)) * camera.aspect), 0.9, 1.6);
    closePos.set(body.x + 0.12 * S, headY + 0.08 * S, body.z + near * S);
    closeLook.set(body.x, headY + 0.01 * S, body.z);
    // 草原場景：鏡頭只往前移一段、其餘用望遠拉近，背景才不會被拉扯變形（遠景略帶模糊，像人像照的景深）
    if (env.wide) { closeLook.set(body.x, headY - 0.15 * S, body.z); closePos.lerpVectors(widePos, closeLook, CLOSE_FRAC); }
    shot = ease(shot, state === 'idle' ? 0 : 1, 1.2, dt);
    camPos.lerpVectors(widePos, closePos, shot);
    camPos.x += Math.sin(t * 0.15) * (env.wide ? 0.03 : 0.15) * (1 - shot * 0.8) * m;   // 草原場景晃動小一點，全景圖才不會變形
    camLook.lerpVectors(wideLook, closeLook, shot);
    camera.position.lerp(camPos, Math.min(1, dt * 3));
    if (env.wide) {
      if (!look.dragging && zoom.target >= 1 && (look.idle += dt) > 3) {   // 放開 3 秒後慢慢轉回來（走最短的方向）；拉近時停在原處
        look.yaw = Math.atan2(Math.sin(look.yaw), Math.cos(look.yaw));
        look.yaw = ease(look.yaw, 0, 1.2, dt); look.pitch = ease(look.pitch, 0, 1.2, dt);
      }
      lookDir.subVectors(camLook, camera.position);
      lookRight.crossVectors(lookDir, camera.up).normalize();
      lookDir.applyAxisAngle(lookRight, look.pitch).applyAxisAngle(camera.up, look.yaw);
      camLook.copy(camera.position).add(lookDir);
    }
    camera.lookAt(camLook);
    zoom.cur = Math.exp(ease(Math.log(zoom.cur), Math.log(zoom.target), 8, dt));   // 平滑過渡
    if (env.wideHFov) {   // 草原場景：遠景放寬視角，構圖和全景圖網站上一樣；回答時用望遠拉近
      const wideV = THREE.MathUtils.clamp(2 * THREE.MathUtils.radToDeg(Math.atan(Math.tan(THREE.MathUtils.degToRad(env.wideHFov / 2)) / camera.aspect)), 35, 60);
      const dist = camera.position.distanceTo(closeLook);
      const closeV = 2 * THREE.MathUtils.radToDeg(Math.atan((camera.aspect < 1 ? 1.5 : 1.1) * S / 2 / dist));   // 框住上半身和捧茶的雙手（手機畫面窄，多留一些）
      const fov = wideV * Math.pow(closeV / wideV, shot) * zoom.cur;   // 以等比例縮放，拉近的速度看起來平均；再乘上滾輪縮放
      if (Math.abs(fov - camera.fov) > 0.01) { camera.fov = fov; camera.updateProjectionMatrix(); }
    }
    else if (Math.abs(baseFov * zoom.cur - camera.fov) > 0.01) { camera.fov = baseFov * zoom.cur; camera.updateProjectionMatrix(); }
    if (debugCam) { camera.position.set(debugCam[0], debugCam[1], debugCam[2]); camera.lookAt(debugCam[3], debugCam[4], debugCam[5]); }
    gaze.position.copy(camera.position);

    vrm.update(dt);
    // 茶杯在兩手之間，跟著手走，杯口始終朝上
    hands[0].getWorldPosition(cupPos).add(hands[1].getWorldPosition(cup.group.position)).multiplyScalar(0.5);
    cup.group.position.copy(cupPos).add(dirBody.copy(CUP_OFFSET).multiplyScalar(S).applyQuaternion(vrm.scene.quaternion));
    cup.group.quaternion.copy(vrm.scene.quaternion);
    cup.update(t, m);
    renderer.render(scene, camera);
  });

  return {
    /** idle | thinking | talking */
    setState(s) { state = s; pendingIdle = false; voiced = false; if (s !== 'talking') queue.length = 0; },
    /** 沒有聲音時：把剛串流到的字幕排進對嘴佇列 */
    speak(text) { state = 'talking'; pendingIdle = false; queue.push(...text); },
    /** 語音開始／唸完一句：唸的時候嘴巴跟著開合 */
    voice(on) { voiced = on; if (on) { state = 'talking'; pendingIdle = false; } },
    /** 回答結束：嘴型停下後回到冥想 */
    finish() { if (queue.length === 0 && !voiced) { state = 'idle'; onIdle?.(); } else pendingIdle = true; },
    /** 停止繪製並釋放資源 */
    dispose() {
      renderer.setAnimationLoop(null);
      resizeObs.disconnect();
      canvas.removeEventListener('pointerdown', onDown);
      canvas.removeEventListener('pointermove', onMove);
      canvas.removeEventListener('pointerup', onUp);
      canvas.removeEventListener('pointercancel', onUp);
      canvas.removeEventListener('wheel', onWheel);
      visibleObs.disconnect();
      VRMUtils.deepDispose(vrm.scene);
      cup.dispose();
      env.dispose();
      scene.traverse(o => { o.geometry?.dispose(); o.material?.dispose?.(); });
      renderer.dispose();
    },
  };
}

// 一只米白瓷茶杯：杯身用旋轉體做出圈足與微微外撇的杯口，杯裡是茶湯，上方幾縷熱氣慢慢升起
function makeTeaCup() {
  const group = new THREE.Group();
  const R = 0.034, H = 0.05;
  const profile = [
    [0, 0], [R * 0.55, 0], [R * 0.55, H * 0.08], [R * 0.62, H * 0.1],
    [R * 0.86, H * 0.35], [R * 0.97, H * 0.75], [R, H],
    [R * 0.93, H], [R * 0.9, H * 0.75], [R * 0.79, H * 0.35], [R * 0.52, H * 0.14], [0, H * 0.14],
  ].map(([x, y]) => new THREE.Vector2(x, y));
  // 夕陽在角色背後，杯子正面偏暗，加一點自發光讓米白釉色看得出來
  const glaze = new THREE.MeshStandardMaterial({ color: 0xf2ead8, emissive: 0xf2ead8, emissiveIntensity: 0.5, roughness: 0.35 });
  const body = new THREE.Mesh(new THREE.LatheGeometry(profile, 32), glaze);
  body.position.y = -H / 2;
  group.add(body);

  const tea = new THREE.Mesh(new THREE.CircleGeometry(R * 0.9, 32),
    new THREE.MeshStandardMaterial({ color: 0xa8743a, emissive: 0xa8743a, emissiveIntensity: 0.25, roughness: 0.15 }));
  tea.rotation.x = -Math.PI / 2;
  tea.position.y = -H / 2 + H * 0.72;
  group.add(tea);

  // 熱氣：柔邊的半透明貼片，各自錯開相位往上飄、淡出
  const c = document.createElement('canvas');
  c.width = c.height = 64;
  const g = c.getContext('2d');
  const grad = g.createRadialGradient(32, 32, 0, 32, 32, 32);
  grad.addColorStop(0, 'rgba(255,255,255,0.55)');
  grad.addColorStop(1, 'rgba(255,255,255,0)');
  g.fillStyle = grad;
  g.fillRect(0, 0, 64, 64);
  const tex = new THREE.CanvasTexture(c);
  const puffs = Array.from({ length: 4 }, (_, i) => {
    const sp = new THREE.Sprite(new THREE.SpriteMaterial({ map: tex, transparent: true, depthWrite: false, opacity: 0 }));
    sp.userData.offset = i / 4;
    group.add(sp);
    return sp;
  });

  return {
    group,
    update(t, m) {
      for (const sp of puffs) {
        const k = (t * 0.22 * (m < 1 ? 0.5 : 1) + sp.userData.offset) % 1;     // 0 → 1：從杯口升到上方
        sp.position.set(Math.sin(t * 0.9 + sp.userData.offset * 6) * 0.01, H / 2 + k * 0.12, 0);
        const size = 0.025 + k * 0.05;
        sp.scale.set(size, size, 1);
        sp.material.opacity = Math.sin(k * Math.PI) * 0.35;
      }
    },
    dispose() { tex.dispose(); },
  };
}
