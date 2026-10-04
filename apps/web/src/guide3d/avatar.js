// 來源：bodhi-guide 首頁的解說員（guide-hero-assets 分支 guide-hero/avatar.js）。
// 官網版多了 dispose()，離開 AI 組長頁時停止繪製並釋放 WebGL 資源。
// 覺行小組線上組長 3D 角色：站在及膝的湖水中冥想。被提問時睜眼、鏡頭推近、依字幕對嘴；
// 回答時在湖中蘆葦區隨機走動，雙手始終捧著一杯熱茶。回答結束後在原地停下、回到冥想。
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { VRMLoaderPlugin, VRMUtils } from '@pixiv/three-vrm';
import { buildScene, WATER_Y, STAND_Z, SHORE_Z, channelHalfWidth } from './scene.js';

const VOWELS = ['aa', 'ih', 'ou', 'ee', 'oh'];
const CHARS_PER_SEC = 7;          // 對嘴速度，之後接 TTS 時改由語音時間軸驅動
const PAUSE = /[\s，。、；：！？,.;:!?「」『』（）()…—\n]/;

// 手臂姿勢（各骨骼在「身體座標」中指向的方向；面向 +Z，角色的左手邊是 +X）：
// 雙手彎肘往前、往中間收，在胸口下方一起捧著熱茶杯，左右對稱
const ARMS = {
  leftUpperArm:  [0.12, -0.95, 0.28],
  rightUpperArm: [-0.12, -0.95, 0.28],
  leftLowerArm:  [-0.35, 0.15, 0.92],
  rightLowerArm: [0.35, 0.15, 0.92],
};
const ARM_TWIST = { rightLowerArm: Math.PI / 2, leftLowerArm: -Math.PI / 2 };   // 沿前臂轉 90°，兩手掌心相對、拇指朝上
const ARM_BONES = Object.keys(ARMS);
// 捧杯：兩手四指繞著杯身彎曲（左手鏡像，彎曲方向相反）
const CURL = { IndexProximal: 0.95, IndexIntermediate: 1.0, IndexDistal: 0.5, MiddleProximal: 1.0, MiddleIntermediate: 1.0, MiddleDistal: 0.5,
  RingProximal: 1.05, RingIntermediate: 1.0, RingDistal: 0.5, LittleProximal: 1.1, LittleIntermediate: 1.0, LittleDistal: 0.5 };
const GRIP = Object.fromEntries(Object.entries(CURL).flatMap(([k, v]) => [[`right${k}`, v], [`left${k}`, -v]]));
// 拇指：繞前臂扭轉後會朝上翹，改成往前彎、貼著杯身外側（繞 Z 軸彎，和四指一樣左手取負）
const THUMBS = [['ThumbMetacarpal', 0.15], ['ThumbProximal', 0.35], ['ThumbDistal', 0.2]];
const CUP_OFFSET = new THREE.Vector3(0, 0.025, 0.05);   // 杯子中心相對兩手腕中點（身體座標）：往前到掌心之間
const restDir = name => new THREE.Vector3(...(name.startsWith('left') ? [1, 0, 0] : [-1, 0, 0]));

// 走動範圍與速度：只在湖中的蘆葦區。上限要讓跟拍鏡頭（角色前方 3.2 公尺）停在岸邊麥田之前
const CLOSE_DIST = 3.2;
const AREA = { xMin: -2.6, xMax: 2.6, zMin: -3.2, zMax: SHORE_Z - 0.05 - CLOSE_DIST };
const WALK_SPEED = 0.36;          // 公尺／秒，在水中慢慢走
const STRIDE = 0.9;               // 一個完整步伐循環走的距離

const smoothstep = (a, b, x) => { const t = Math.min(1, Math.max(0, (x - a) / (b - a))); return t * t * (3 - 2 * t); };
const angleLerp = (a, b, k) => { let d = ((b - a + Math.PI) % (Math.PI * 2)) - Math.PI; if (d < -Math.PI) d += Math.PI * 2; return a + d * k; };

export async function createAvatar(canvas, url, { onProgress, onIdle } = {}) {
  const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  renderer.outputColorSpace = THREE.SRGBColorSpace;

  const scene = new THREE.Scene();
  const camera = new THREE.PerspectiveCamera(35, 1, 0.05, 1000);
  const env = buildScene(scene);

  const loader = new GLTFLoader();
  loader.register(parser => new VRMLoaderPlugin(parser));
  const gltf = await loader.loadAsync(url, e => e.total && onProgress?.(e.loaded / e.total));
  const vrm = gltf.userData.vrm;
  VRMUtils.removeUnnecessaryVertices(gltf.scene);
  VRMUtils.rotateVRM0(vrm);
  vrm.scene.traverse(o => { o.frustumCulled = false; });
  scene.add(vrm.scene);

  const bone = name => vrm.humanoid.getNormalizedBoneNode(name);
  const B = Object.fromEntries(['hips', 'spine', 'chest', 'neck', 'head',
    'leftUpperLeg', 'rightUpperLeg', 'leftLowerLeg', 'rightLowerLeg', ...ARM_BONES].map(n => [n, bone(n)]));
  for (const [name, curl] of Object.entries(GRIP)) bone(name)?.rotation.set(0, 0, curl);
  // 拇指扣在杯緣
  for (const [name, z] of THUMBS) { bone(`right${name}`)?.rotation.set(0, 0, z); bone(`left${name}`)?.rotation.set(0, 0, -z); }
  const rawHands = ['rightHand', 'leftHand'].map(n => vrm.humanoid.getRawBoneNode(n));

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
  cup.group.scale.setScalar(1.25);   // 雙手捧著，杯子比單手握時大一點
  scene.add(cup.group);
  const cupPos = new THREE.Vector3();

  // 站進水裡：讓水面剛好在膝蓋（小腿骨的起點）高度；走上岸時沿著坡度出水
  vrm.scene.position.set(0, 0, STAND_Z);
  vrm.scene.updateMatrixWorld(true);
  const kneeY = B.leftLowerLeg.getWorldPosition(new THREE.Vector3()).y;
  const depthAt = z => kneeY * (1 - smoothstep(SHORE_Z - 0.3, SHORE_Z + 0.2, z));
  vrm.scene.position.y = WATER_Y - depthAt(STAND_Z);
  vrm.update(0);
  vrm.scene.updateMatrixWorld(true);

  // 蘆葦最高只到肩膀（上臂骨的起點）
  env.addVegetation(B.leftUpperArm.getWorldPosition(new THREE.Vector3()).y);
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

  let state = 'idle', pendingIdle = false;
  const queue = [];
  const mouth = { vowel: 'aa', amount: 0, target: 0 };
  const face = { eyesClosed: 1, happy: 0, relaxed: 0 };
  let nextBlink = 3, blinkT = -1, shot = 0;   // shot: 0 = 遠景，1 = 近景

  // 走動狀態
  const body = { x: 0, z: STAND_Z, yaw: 0, walk: 0, phase: 0 };
  let waypoint = null, pauseLeft = 0;
  const pickWaypoint = () => {
    for (let i = 0; i < 20; i++) {
      // 在湖中走到視線通道邊緣附近，身旁就是蘆葦叢
      const z = AREA.zMin + Math.random() * (AREA.zMax - AREA.zMin);
      let x = Math.sign(Math.random() - 0.5) * Math.max(0.3, channelHalfWidth(z) - 0.2 - Math.random() * 0.6);
      x = THREE.MathUtils.clamp(x, AREA.xMin, AREA.xMax);
      if (Math.hypot(x - body.x, z - body.z) > 1.2) return { x, z };
    }
    return { x: 0, z: STAND_Z };
  };

  const resize = () => {
    const w = canvas.clientWidth, h = canvas.clientHeight;
    if (!w || !h) return;
    renderer.setSize(w, h, false);
    camera.aspect = w / h;
    // 鏡頭平移（不傾斜）：把畫面往上挪，角色落在中上方，下方留給字幕與輸入框
    camera.setViewOffset(w, h, 0, Math.round(h * 0.14), w, h);
    camera.updateProjectionMatrix();
  };
  const resizeObs = new ResizeObserver(resize);
  resizeObs.observe(canvas);
  resize();

  // 捲到頁面下方、看不到場景時暫停繪製
  let onScreen = true;
  const visibleObs = new IntersectionObserver(([e]) => { onScreen = e.isIntersecting; });
  visibleObs.observe(canvas);

  const camPos = new THREE.Vector3(), camLook = new THREE.Vector3(), wideLook = new THREE.Vector3(), closeLook = new THREE.Vector3();
  const widePos = new THREE.Vector3(), closePos = new THREE.Vector3(), dirBody = new THREE.Vector3(), hipsW = new THREE.Vector3();
  // 除錯用：?cam=x,y,z,lookX,lookY,lookZ 固定鏡頭
  const debugCam = new URLSearchParams(location.search).get('cam')?.split(',').map(Number);
  if (debugCam || new URLSearchParams(location.search).has('debug')) window.__bodhi = { vrm, THREE, body, camera };
  camera.position.set(0.5, vrm.scene.position.y + headOffset + 0.25, 6.4);   // 從遠景開始，避免第一幀從原點飛進來
  const clock = new THREE.Clock();
  let t = 0, charClock = 0;

  renderer.setAnimationLoop(() => {
    const dt = Math.min(clock.getDelta(), 0.1);
    if (!onScreen) return;
    t += dt;
    env.uTime.value = t;
    const m = reduceMotion ? 0.3 : 1;

    // ---- 走動：回答時在蘆葦與麥田間隨機走，停下時轉身面向觀眾 ----
    let moving = false, faceYaw;
    if (state === 'talking') {
      if (!waypoint) {
        pauseLeft -= dt;
        if (pauseLeft <= 0) waypoint = pickWaypoint();
      }
      if (waypoint) {
        const dx = waypoint.x - body.x, dz = waypoint.z - body.z, dist = Math.hypot(dx, dz);
        if (dist < 0.08) { waypoint = null; pauseLeft = 0.8 + Math.random() * 1.6; }
        else {
          moving = true;
          const step = Math.min(dist, WALK_SPEED * body.walk * dt * m);
          body.x += dx / dist * step; body.z += dz / dist * step;
          faceYaw = Math.atan2(dx, dz);
        }
      }
    } else {
      waypoint = null; pauseLeft = 0;
    }
    if (faceYaw === undefined) faceYaw = Math.atan2(camera.position.x - body.x, camera.position.z - body.z);
    body.walk = ease(body.walk, moving ? 1 : 0, 3, dt);
    body.yaw = angleLerp(body.yaw, faceYaw, Math.min(1, 2.5 * dt));
    body.phase += dt * body.walk * (Math.PI * 2) * (WALK_SPEED / STRIDE) * m;

    const depth = depthAt(body.z);
    vrm.scene.position.set(body.x, WATER_Y - depth + Math.abs(Math.sin(body.phase)) * 0.018 * body.walk, body.z);
    vrm.scene.rotation.y = body.yaw;
    vrm.scene.updateMatrixWorld(true);
    qBodyInv.copy(vrm.scene.quaternion).invert();

    // ---- 腿：走路的擺動（在水裡步伐小一點）----
    const sw = Math.sin(body.phase) * 0.34 * body.walk;
    if (B.leftUpperLeg) B.leftUpperLeg.rotation.set(-sw, 0, 0);
    if (B.rightUpperLeg) B.rightUpperLeg.rotation.set(sw, 0, 0);
    if (B.leftLowerLeg) B.leftLowerLeg.rotation.set(Math.max(0, -Math.cos(body.phase)) * 0.55 * body.walk, 0, 0);
    if (B.rightLowerLeg) B.rightLowerLeg.rotation.set(Math.max(0, Math.cos(body.phase)) * 0.55 * body.walk, 0, 0);
    if (B.hips) B.hips.rotation.y = Math.sin(body.phase) * 0.06 * body.walk;

    // ---- 身體、頭、手臂 ----
    const breath = Math.sin(t * (state === 'idle' ? 1.1 : 1.6));
    if (B.chest) B.chest.rotation.x = -0.03 + breath * 0.02 * m;
    if (B.spine) B.spine.rotation.set(0.04 + body.walk * 0.05, -Math.sin(body.phase) * 0.05 * body.walk, 0);
    vrm.scene.updateMatrixWorld(true);
    for (const name of ARM_BONES) aimBody(name, dirBody.set(...ARMS[name]).normalize(), ARM_TWIST[name]);
    let hx = state === 'idle' ? 0.12 : 0.02, hy = 0, hz = 0;
    if (state === 'thinking') { hz = 0.08; hx = 0.05; }
    if (state === 'talking') { hx += Math.sin(t * 2.3) * 0.025 * m; hy = Math.sin(t * 1.3) * 0.04 * m; }
    B.head.rotation.x = ease(B.head.rotation.x, hx, 3, dt);
    B.head.rotation.y = ease(B.head.rotation.y, hy, 3, dt);
    B.head.rotation.z = ease(B.head.rotation.z, hz, 3, dt);

    // ---- 身體走過的地方，蘆葦往兩旁分開 ----
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

    // ---- 對嘴 ----
    charClock += dt;
    if (charClock >= 1 / CHARS_PER_SEC) {
      charClock = 0;
      const ch = queue.shift();
      if (ch === undefined || PAUSE.test(ch)) mouth.target = 0;
      else { mouth.vowel = VOWELS[ch.codePointAt(0) % VOWELS.length]; mouth.target = 0.45 + (ch.codePointAt(0) % 7) / 14; }
    }
    mouth.amount = ease(mouth.amount, mouth.target, 18, dt);
    for (const v of VOWELS) set(v, v === mouth.vowel ? mouth.amount : ease(em?.getValue(v) ?? 0, 0, 18, dt));
    if (state === 'talking' && queue.length === 0 && mouth.amount < 0.02 && pendingIdle) { state = 'idle'; pendingIdle = false; onIdle?.(); }

    // ---- 風吹頭髮 ----
    const gust = 0.06 + 0.05 * (Math.sin(t * 1.3) * 0.6 + Math.sin(t * 2.6) * 0.25 + 0.35) * m;
    for (const j of joints) { j.settings.gravityDir.copy(windDir); j.settings.gravityPower = gust; }

    // ---- 鏡頭：跟著角色；回答時推近、冥想時拉遠，並緩慢漂移 ----
    const headY = vrm.scene.position.y + headOffset;
    widePos.set(0.5 + body.x * 0.35, headY + 0.25, 6.4);
    wideLook.set(body.x, headY - 0.1, body.z - 1.5);
    // 近景保持約 3.2 公尺：看得到全身動作與身旁的蘆葦，不貼臉
    closePos.set(body.x + 0.5, headY + 0.05, body.z + CLOSE_DIST);
    closeLook.set(body.x, headY - 0.38, body.z);
    shot = ease(shot, state === 'idle' ? 0 : 1, 1.2, dt);
    camPos.lerpVectors(widePos, closePos, shot);
    camPos.x += Math.sin(t * 0.15) * 0.15 * m;
    camLook.lerpVectors(wideLook, closeLook, shot);
    camera.position.lerp(camPos, Math.min(1, dt * 3));
    camera.lookAt(camLook);
    if (debugCam) { camera.position.set(debugCam[0], debugCam[1], debugCam[2]); camera.lookAt(debugCam[3], debugCam[4], debugCam[5]); }
    gaze.position.copy(camera.position);

    vrm.update(dt);
    // 茶杯在兩手之間，杯口始終朝上
    rawHands[0].getWorldPosition(cupPos).add(rawHands[1].getWorldPosition(dirBody)).multiplyScalar(0.5);
    cup.group.position.copy(cupPos).add(dirBody.copy(CUP_OFFSET).applyQuaternion(vrm.scene.quaternion));
    cup.group.quaternion.copy(vrm.scene.quaternion);
    cup.update(t, m);
    renderer.render(scene, camera);
  });

  return {
    /** idle | thinking | talking */
    setState(s) { state = s; pendingIdle = false; if (s !== 'talking') queue.length = 0; },
    /** 把剛串流到的字幕排進對嘴佇列 */
    speak(text) { state = 'talking'; pendingIdle = false; queue.push(...text); },
    /** 回答結束：唸完佇列裡的字後，在原地停下回到冥想 */
    finish() { if (queue.length === 0) { state = 'idle'; onIdle?.(); } else pendingIdle = true; },
    /** 停止繪製並釋放資源 */
    dispose() {
      renderer.setAnimationLoop(null);
      resizeObs.disconnect();
      visibleObs.disconnect();
      VRMUtils.deepDispose(vrm.scene);
      cup.dispose();
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
