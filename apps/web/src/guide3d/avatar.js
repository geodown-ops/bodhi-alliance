// 來源：bodhi-guide 首頁的解說員（guide-hero-assets 分支 guide-hero/avatar.js）。
// 官網版多了 dispose()，離開 AI 組長頁時停止繪製並釋放 WebGL 資源。
// 覺行小組線上組長 3D 角色：站在及膝的湖水中冥想。被提問時睜眼、鏡頭推近、依字幕對嘴；
// 回答時在湖中蘆葦區隨機走動，右手始終端著一杯茶、左手自然垂下。回答結束後在原地停下、回到冥想。
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { VRMLoaderPlugin, VRMUtils } from '@pixiv/three-vrm';
import { buildScene, WATER_Y, STAND_Z, SHORE_Z, channelHalfWidth } from './scene.js';

const VOWELS = ['aa', 'ih', 'ou', 'ee', 'oh'];
const CHARS_PER_SEC = 7;          // 對嘴速度，之後接 TTS 時改由語音時間軸驅動
const PAUSE = /[\s，。、；：！？,.;:!?「」『』（）()…—\n]/;

// 手臂姿勢：左手自然垂在身側，右手彎肘在胸前端著一杯茶（各骨骼在「身體座標」中指向的方向；面向 +Z，角色的左手邊是 +X）
const ARMS = {
  leftUpperArm:  [0.24, -0.97, 0.02],     // 略往外張，避開寬鬆上衣
  rightUpperArm: [-0.2, -0.95, 0.2],
  leftLowerArm:  [0.2, -0.96, 0.15],      // 手掌落在大腿外側，不穿進短褲
  rightLowerArm: [0.3, 0.12, 0.95],       // 前臂往前、略往身體中線，茶杯在胸口下方
};
const ARM_BONES = Object.keys(ARMS);
const CUP = { along: 0.06, inward: 0.035, up: 0.0 };   // 茶杯中心相對右手腕：沿前臂、往身體中線、往上（公尺）
const RIGHT_FINGERS = ['Index', 'Middle', 'Ring', 'Little'].flatMap(f => ['Proximal', 'Intermediate', 'Distal'].map(s => `right${f}${s}`));

// 茶杯：無把手的陶瓷茶杯，杯裡是茶湯。全部程序化，不用貼圖
function makeTeaCup() {
  const cup = new THREE.Group();
  const H = 0.06, RT = 0.03, RB = 0.022;
  const glaze = new THREE.MeshStandardMaterial({ color: '#f2ead8', emissive: '#f2ead8', emissiveIntensity: 0.45, roughness: 0.35 });   // 夕陽在背後，自發光補一點亮度
  const wall = new THREE.Mesh(new THREE.CylinderGeometry(RT, RB, H, 28, 1, true), glaze);
  wall.material.side = THREE.DoubleSide;
  const foot = new THREE.Mesh(new THREE.CylinderGeometry(RB, RB * 0.9, 0.006, 28), glaze);
  foot.position.y = -H / 2;
  const rim = new THREE.Mesh(new THREE.TorusGeometry(RT, 0.0025, 8, 32), glaze);
  rim.rotation.x = Math.PI / 2; rim.position.y = H / 2;
  const tea = new THREE.Mesh(new THREE.CircleGeometry(RT * 0.94, 28), new THREE.MeshStandardMaterial({ color: '#a8743a', emissive: '#a8743a', emissiveIntensity: 0.3, roughness: 0.15 }));
  tea.rotation.x = -Math.PI / 2; tea.position.y = H / 2 - 0.012;
  cup.add(wall, foot, rim, tea);
  return cup;
}
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
    'leftUpperLeg', 'rightUpperLeg', 'leftLowerLeg', 'rightLowerLeg', 'rightHand', ...ARM_BONES].map(n => [n, bone(n)]));

  // 右手握杯：手掌轉向身體中線、四指彎曲圍住杯身
  if (B.rightHand) B.rightHand.rotation.x = -Math.PI / 2;
  for (const n of RIGHT_FINGERS) { const f = bone(n); if (f) f.rotation.z = n.endsWith('Proximal') ? 0.9 : 0.7; }
  for (const [n, x] of [['rightThumbProximal', 1.0], ['rightThumbIntermediate', 0.4]]) { const f = bone(n); if (f) f.rotation.x = x; }   // 拇指扣在杯緣外側
  const cup = makeTeaCup();
  scene.add(cup);
  const handW = new THREE.Vector3(), forearmW = new THREE.Vector3(), inwardW = new THREE.Vector3();

  // 依「身體座標」中的目標方向擺手臂：扣掉整個角色的轉向，求出骨骼的區域旋轉
  const qBodyInv = new THREE.Quaternion(), qParent = new THREE.Quaternion(), qTarget = new THREE.Quaternion();
  const aimBody = (name, dirBody) => {
    const b = B[name]; if (!b) return;
    b.parent.getWorldQuaternion(qParent).premultiply(qBodyInv);
    qTarget.setFromUnitVectors(restDir(name), dirBody);
    b.quaternion.copy(qParent.invert().multiply(qTarget));
    b.updateMatrixWorld(true);
  };

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
    for (const name of ARM_BONES) aimBody(name, dirBody.set(...ARMS[name]).normalize());
    let hx = state === 'idle' ? 0.12 : 0.02, hy = 0, hz = 0;
    if (state === 'thinking') { hz = 0.08; hx = 0.05; }
    if (state === 'talking') { hx += Math.sin(t * 2.3) * 0.025 * m; hy = Math.sin(t * 1.3) * 0.04 * m; }
    B.head.rotation.x = ease(B.head.rotation.x, hx, 3, dt);
    B.head.rotation.y = ease(B.head.rotation.y, hy, 3, dt);
    B.head.rotation.z = ease(B.head.rotation.z, hz, 3, dt);

    // ---- 身體走過的地方，蘆葦往兩旁分開 ----
    vrm.scene.updateMatrixWorld(true);

    // ---- 茶杯跟著右手，但永遠保持直立、不灑出來 ----
    if (B.rightHand) {
      B.rightHand.getWorldPosition(handW);
      forearmW.set(...ARMS.rightLowerArm).normalize().applyQuaternion(vrm.scene.quaternion);
      inwardW.set(1, 0, 0).applyQuaternion(vrm.scene.quaternion);
      cup.position.copy(handW).addScaledVector(forearmW, CUP.along).addScaledVector(inwardW, CUP.inward);
      cup.position.y += CUP.up;
      cup.rotation.set(0, body.yaw, 0);
    }
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
      scene.traverse(o => { o.geometry?.dispose(); o.material?.dispose?.(); });
      renderer.dispose();
    },
  };
}
