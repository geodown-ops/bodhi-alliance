// 來源：bodhi-guide 首頁的解說員（guide-hero-assets 分支 guide-hero/scene.js）。
// 冥想空間的黃昏版（下午五點到隔天早上六點）：黃昏、湖面、金色緩丘，水中與岸邊長滿隨風擺動的白色蘆葦。
// 白天的藍天湖景在 scene.js；兩者提供同一組介面，avatar.js 依時間挑一個。
// 全部以程序化幾何與著色器生成，不使用外部貼圖。
import * as THREE from 'three';
import { WATER_Y, STAND_Z, SHORE_Z, channelHalfWidth } from './scene.js';

const SKY_TOP = new THREE.Color('#9fbfd6');
const SKY_HORIZON = new THREE.Color('#f4d3a6');
const SUN_COLOR = new THREE.Color('#ffd9a0');
const FOG_COLOR = new THREE.Color('#ead0ab');

// 會把植物撥開的點：覺行小組線上組長的雙手與身體（xyz = 位置，w = 影響半徑；w = 0 表示不作用）
const PUSH = { value: [new THREE.Vector4(), new THREE.Vector4(), new THREE.Vector4()] };

export function buildDuskScene(scene) {
  const uTime = { value: 0 };
  const sunDir = new THREE.Vector3(-0.55, 0.16, -0.82).normalize();   // 低角度的夕陽，從左後方

  scene.fog = new THREE.Fog(FOG_COLOR, 35, 190);

  // 天空：上淡藍、地平線暖橘，太陽方向一圈光暈
  const sky = new THREE.Mesh(
    new THREE.SphereGeometry(400, 32, 16),
    new THREE.ShaderMaterial({
      side: THREE.BackSide, depthWrite: false, fog: false,
      uniforms: { top: { value: SKY_TOP }, horizon: { value: SKY_HORIZON }, sun: { value: SUN_COLOR }, sunDir: { value: sunDir } },
      vertexShader: `varying vec3 vDir; void main(){ vDir = normalize(position); gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.); }`,
      fragmentShader: `uniform vec3 top, horizon, sun, sunDir; varying vec3 vDir;
        void main(){
          float h = clamp(vDir.y, 0., 1.);
          vec3 c = mix(horizon, top, pow(h, .55));
          float s = max(dot(normalize(vDir), sunDir), 0.);
          c += sun * (pow(s, 64.) * .9 + pow(s, 6.) * .18);
          gl_FragColor = vec4(c, 1.);
        }`,
    }));
  scene.add(sky);

  // 遠方緩丘
  const hillMat = new THREE.MeshLambertMaterial({ color: '#c79c63' });
  for (const [x, z, sx, sy, sz] of [[-70, -150, 110, 14, 40], [40, -170, 120, 18, 45], [120, -140, 70, 10, 35], [-140, -120, 80, 9, 30], [0, -210, 160, 26, 50]]) {
    const hill = new THREE.Mesh(new THREE.SphereGeometry(1, 48, 24), hillMat);
    hill.scale.set(sx, sy, sz); hill.position.set(x, -1, z);
    scene.add(hill);
  }

  // 湖面：天空倒影 + 菲涅耳 + 細碎波光
  const water = new THREE.Mesh(
    new THREE.PlaneGeometry(600, 600, 1, 1).rotateX(-Math.PI / 2),
    new THREE.ShaderMaterial({
      fog: false,
      uniforms: { uTime, top: { value: SKY_TOP }, horizon: { value: SKY_HORIZON }, sun: { value: SUN_COLOR }, sunDir: { value: sunDir }, fogColor: { value: FOG_COLOR } },
      vertexShader: `varying vec3 vWorld; void main(){ vec4 w = modelMatrix * vec4(position,1.); vWorld = w.xyz; gl_Position = projectionMatrix * viewMatrix * w; }`,
      fragmentShader: `uniform float uTime; uniform vec3 top, horizon, sun, sunDir, fogColor; varying vec3 vWorld;
        void main(){
          vec2 p = vWorld.xz;
          float r = sin(p.x*1.7 + uTime*.9) * .5 + sin(p.y*2.3 - uTime*1.1) * .35 + sin((p.x+p.y)*3.1 + uTime*1.7) * .15;
          vec3 n = normalize(vec3(r*.04, 1., r*.03));
          vec3 v = normalize(cameraPosition - vWorld);
          vec3 refl = reflect(-v, n);
          float fres = pow(1. - max(dot(n, v), 0.), 3.);
          vec3 skyCol = mix(horizon, top, pow(clamp(refl.y,0.,1.), .55));
          vec3 deep = vec3(.42, .52, .58);
          vec3 c = mix(deep, skyCol, .55 + .45*fres);
          c += sun * pow(max(dot(refl, sunDir), 0.), 180.) * 1.4;
          float d = length(cameraPosition.xz - vWorld.xz);
          c = mix(c, fogColor, smoothstep(40., 200., d));
          gl_FragColor = vec4(c, 1.);
        }`,
    }));
  scene.add(water);

  // 角色腳邊的漣漪
  const ripple = ripples(uTime);
  scene.add(ripple);

  // 燈光：夕陽逆光 + 天光 + 正面補光（避免臉太暗）
  scene.add(new THREE.HemisphereLight('#c6d8e6', '#b08a52', 1.1));
  const sunLight = new THREE.DirectionalLight('#ffcf94', 2.1);
  sunLight.position.copy(sunDir).multiplyScalar(20);
  scene.add(sunLight);
  const fill = new THREE.DirectionalLight('#fff1dd', 0.9);
  fill.position.set(1.5, 2.5, 5);
  scene.add(fill);

  return {
    uTime, sunDir,
    /**
     * 植物要等角色載入後才知道肩膀高度：白色蘆葦最高不超過肩膀，最前景補一條岸邊麥田
     * @param {number} maxReedTop 蘆葦（含花穗）的最高點，世界座標 y
     */
    addVegetation(maxReedTop) {
      for (const mesh of reeds(uTime, maxReedTop)) scene.add(mesh);
      for (const mesh of wheat(uTime)) scene.add(mesh);
    },
    /** 漣漪跟著角色移動；走上岸時淡出 */
    setRipple(x, z, strength) {
      ripple.position.x = x; ripple.position.z = z;
      ripple.material.uniforms.uStrength.value = strength;
      ripple.visible = strength > 0.01;
    },
    /** 設定撥動植物的點：[{ position: Vector3, radius }]，最多三個 */
    setPushers(list) {
      PUSH.value.forEach((v, i) => {
        const p = list[i];
        if (p) v.set(p.position.x, p.position.y, p.position.z, p.radius); else v.set(0, 0, 0, 0);
      });
    },
    // 黃昏湖景沒有倒影與動物，下面這幾個是為了和白天版介面一致
    update() {},
    resize() {},
    summon() {},
    dispose() {},
  };
}

// 以世界座標的高度計算風吹位移：莖與穗共用同一個公式，穗才會跟著莖一起擺。
// 另外依 PUSH 的點把附近的植物往外撥開（越高處位移越大，根部不動）。
function windShader(mat, uTime, strength, baseY, plantHeight) {
  mat.onBeforeCompile = s => {
    s.uniforms.uTime = uTime;
    s.uniforms.uPush = PUSH;
    s.vertexShader = 'uniform float uTime;\nuniform vec4 uPush[3];\n' + s.vertexShader.replace('#include <project_vertex>', `
      vec4 wp = modelMatrix * instanceMatrix * vec4(transformed, 1.0);
      vec3 base = vec3(instanceMatrix[3]);
      float hh = max(wp.y - (${baseY.toFixed(3)}), 0.);
      float gust = sin(uTime*1.1 + base.x*.35 + base.z*.25) * .6 + sin(uTime*2.3 + base.x*1.3 + base.z*.7) * .25 + .45;
      wp.x += gust * ${strength.toFixed(3)} * hh * hh;
      wp.z += gust * ${(strength * 0.35).toFixed(3)} * hh * hh;
      float bend = clamp(hh / ${plantHeight.toFixed(3)}, 0., 1.);
      for (int i = 0; i < 3; i++) {
        float r = uPush[i].w;
        if (r <= 0.) continue;
        vec2 dv = base.xz - uPush[i].xz;
        float dl = length(dv);
        float f = (1. - smoothstep(0., r, dl)) * bend * bend;
        wp.xz += (dl > 1e-4 ? dv / dl : vec2(1., 0.)) * f * r * .9;
        wp.y -= f * r * .25;
      }
      vec4 mvPosition = viewMatrix * wp;
      gl_Position = projectionMatrix * mvPosition;`);
  };
}

// 白色蘆葦：一叢一叢分布在水中，避開角色與鏡頭之間的視線；高度介於肩膀的六成到肩膀之間
const PLUME_SCALE = 0.6;                       // 花穗長度比例
const PLUME_TOP = 0.46 * PLUME_SCALE;          // 花穗頂端高出莖頂的距離（球體 0.26 縮放、往上平移 0.2）
function reeds(uTime, maxTop) {
  const baseY = WATER_Y - 0.3;
  let seed = 11; const rnd = () => (seed = (seed * 16807) % 2147483647) / 2147483647;
  const items = [];
  const clumps = 220;
  for (let c = 0; c < clumps; c++) {
    const cx = (rnd() - 0.5) * 36, cz = -18 + rnd() * 20.4;     // 到岸邊麥田前為止
    // 鏡頭到角色之間留一條往前張開的視線通道；角色身後的蘆葦則可以靠近一點，當作背景
    if (Math.abs(cx - 0.15) < channelHalfWidth(cz) && cz > STAND_Z - 1.2) continue;
    const n = 8 + Math.floor(rnd() * 14);
    for (let i = 0; i < n; i++) {
      const a = rnd() * Math.PI * 2, r = Math.sqrt(rnd()) * 0.55;
      items.push({
        x: cx + Math.cos(a) * r, z: cz + Math.sin(a) * r,
        y: baseY,
        h: maxTop * (0.6 + rnd() * 0.4) - baseY - PLUME_TOP,
        tilt: (rnd() - 0.5) * 0.18, yaw: rnd() * Math.PI * 2, droop: 0.25 + rnd() * 0.35,
        tone: 0.92 + rnd() * 0.08,
      });
    }
  }

  const stemGeo = new THREE.PlaneGeometry(0.014, 1, 1, 6).translate(0, 0.5, 0);
  const stemMat = new THREE.MeshLambertMaterial({ color: '#cdb98f', side: THREE.DoubleSide });
  windShader(stemMat, uTime, 0.05, baseY, maxTop - baseY);
  const stems = new THREE.InstancedMesh(stemGeo, stemMat, items.length);

  // 花穗：細長、微微下垂的羽狀穗
  const plumeGeo = new THREE.SphereGeometry(1, 10, 8).scale(0.026, 0.26, 0.02).translate(0, 0.2, 0);
  // 逆光下仍要看起來是白色：用自發光補亮，帶一點夕陽的暖色
  const plumeMat = new THREE.MeshLambertMaterial({ color: '#ffffff', emissive: '#efe2cf', emissiveIntensity: 0.62 });
  windShader(plumeMat, uTime, 0.05, baseY, maxTop - baseY);
  const plumes = new THREE.InstancedMesh(plumeGeo, plumeMat, items.length);

  const m = new THREE.Matrix4(), q = new THREE.Quaternion(), s = new THREE.Vector3(1, 1, 1), p = new THREE.Vector3(), e = new THREE.Euler();
  const color = new THREE.Color();
  items.forEach((it, i) => {
    e.set(it.tilt, it.yaw, it.tilt * 0.6);
    q.setFromEuler(e);
    s.set(1, it.h, 1);
    m.compose(p.set(it.x, it.y, it.z), q, s);
    stems.setMatrixAt(i, m);
    // 花穗接在莖的頂端，往一側下垂
    const top = new THREE.Vector3(0, it.h, 0).applyQuaternion(q).add(p.set(it.x, it.y, it.z));
    e.set(it.tilt + it.droop, it.yaw, it.tilt * 0.6);
    q.setFromEuler(e);
    m.compose(top, q, s.set(1, PLUME_SCALE, 1));
    plumes.setMatrixAt(i, m);
    plumes.setColorAt(i, color.setScalar(it.tone));
  });
  stems.frustumCulled = plumes.frustumCulled = false;
  return [stems, plumes];
}

// 最前景的岸邊麥田：只鋪在鏡頭正前方的一條帶狀區域，高度壓低，框住畫面下緣而不擋到人物
function wheat(uTime) {
  let seed = 23; const rnd = () => (seed = (seed * 16807) % 2147483647) / 2147483647;
  const group = [];
  const groundY = 0.02;

  const bank = new THREE.Mesh(new THREE.PlaneGeometry(40, 10).rotateX(-Math.PI / 2), new THREE.MeshLambertMaterial({ color: '#a8844c' }));
  bank.position.set(0, groundY - 0.01, SHORE_Z + 5);   // 岸邊從湖岸線開始
  group.push(bank);

  const count = 9000;
  const items = [];
  for (let i = 0; i < count; i++) {
    const z = SHORE_Z + 0.15 + Math.pow(rnd(), 0.7) * 1.1;
    items.push({ x: -4.5 + rnd() * 10, z, h: 0.26 + rnd() * 0.2, tilt: (rnd() - 0.5) * 0.3, yaw: rnd() * Math.PI, tone: 0.82 + rnd() * 0.3 });
  }

  const stalkGeo = new THREE.PlaneGeometry(0.005, 1, 1, 4).translate(0, 0.5, 0);
  const stalkMat = new THREE.MeshLambertMaterial({ color: '#d8b268', side: THREE.DoubleSide });
  windShader(stalkMat, uTime, 0.18, groundY, 0.55);
  const stalks = new THREE.InstancedMesh(stalkGeo, stalkMat, count);

  // 麥穗：飽滿的金色穗頭
  const earGeo = new THREE.SphereGeometry(1, 8, 6).scale(0.007, 0.038, 0.007).translate(0, 0.032, 0);
  const earMat = new THREE.MeshLambertMaterial({ color: '#f0c972', emissive: '#6e4a1c', emissiveIntensity: 0.35 });
  windShader(earMat, uTime, 0.18, groundY, 0.55);
  const ears = new THREE.InstancedMesh(earGeo, earMat, count);

  const m = new THREE.Matrix4(), q = new THREE.Quaternion(), s = new THREE.Vector3(), p = new THREE.Vector3(), e = new THREE.Euler();
  const color = new THREE.Color();
  items.forEach((it, i) => {
    q.setFromEuler(e.set(it.tilt, it.yaw, it.tilt * 0.5));
    p.set(it.x, groundY, it.z);
    stalks.setMatrixAt(i, m.compose(p, q, s.set(1, it.h, 1)));
    stalks.setColorAt(i, color.setScalar(it.tone));
    const top = new THREE.Vector3(0, it.h, 0).applyQuaternion(q).add(p);
    q.setFromEuler(e.set(it.tilt + 0.35, it.yaw, it.tilt * 0.5));
    ears.setMatrixAt(i, m.compose(top, q, s.set(1, 1, 1)));
    ears.setColorAt(i, color.setScalar(it.tone));
  });
  stalks.frustumCulled = ears.frustumCulled = false;
  group.push(stalks, ears);
  return group;
}

// 腳邊緩緩擴散的漣漪
function ripples(uTime) {
  const mesh = new THREE.Mesh(
    new THREE.PlaneGeometry(2.4, 2.4).rotateX(-Math.PI / 2),
    new THREE.ShaderMaterial({
      transparent: true, depthWrite: false,
      uniforms: { uTime, uStrength: { value: 1 } },
      vertexShader: `varying vec2 vUv; void main(){ vUv = uv; gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.); }`,
      fragmentShader: `uniform float uTime, uStrength; varying vec2 vUv;
        void main(){
          float d = length((vUv - .5) * 2.);
          float a = 0.;
          for (int i = 0; i < 3; i++) {
            float r = fract(uTime * .18 + float(i) / 3.);
            a += smoothstep(.035, 0., abs(d - (.18 + r * .82))) * (1. - r);
          }
          a *= smoothstep(1., .7, d) * smoothstep(.1, .2, d);
          gl_FragColor = vec4(vec3(1., .96, .9), a * .35 * uStrength);
        }`,
    }));
  mesh.position.set(0, WATER_Y + 0.004, STAND_Z);
  return mesh;
}
