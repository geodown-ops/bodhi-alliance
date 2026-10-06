// 來源：bodhi-guide 首頁的解說員（guide-hero-assets 分支 guide-hero/scene.js），官網版改成白天的湖景。
// 冥想空間：藍天白雲、清澈的湖、綠色緩丘，水中與岸邊長滿隨風擺動的綠色蘆葦與青草。
// 湖水透明看得到湖底光紋與游魚，水面淡淡映出天空與覺行小組線上組長；偶爾有白天鵝游過、燕子成群飛過。
// 全部以程序化幾何與著色器生成，不使用外部貼圖。
import * as THREE from 'three';
import { Reflector } from 'three/addons/objects/Reflector.js';

export const WATER_Y = 0;              // 水面高度
export const STAND_Z = -1.0;           // 角色站在水中的位置
export const SHORE_Z = 3.85;           // 湖岸線：z 大於此處是岸邊草地
/** 會映在水面上的物件要打開這一層（天空、遠山、角色、天鵝、燈光） */
export const REFLECT_LAYER = 1;

/** 鏡頭到角色之間的視線通道半寬：蘆葦不長在通道裡，覺行小組線上組長走到通道邊緣就能撥到蘆葦 */
export function channelHalfWidth(z) {
  return z > STAND_Z ? 1.0 + (z - STAND_Z) * 0.55 : 0.9;
}

const SKY_TOP = new THREE.Color('#2f74c8');
const SKY_HORIZON = new THREE.Color('#a9cfec');
const SUN_COLOR = new THREE.Color('#fff4dc');
const FOG_COLOR = new THREE.Color('#c4dbea');
const WATER_DEEP = new THREE.Color('#1d6f7c');
const WIND = new THREE.Vector2(0.9, 0.35).normalize();   // 微風吹的方向（世界 xz），雲、水波、植物都順著它

// 會把植物撥開的點：0 號是覺行小組線上組長的身體，1、2 號是游過的天鵝（xyz = 位置，w = 影響半徑；w = 0 表示不作用）
const PUSH = { value: [new THREE.Vector4(), new THREE.Vector4(), new THREE.Vector4()] };

const NOISE_GLSL = `
  float hash(vec2 p){ return fract(sin(dot(p, vec2(127.1, 311.7))) * 43758.5453); }
  float noise(vec2 p){ vec2 i = floor(p), f = fract(p); f = f*f*(3.-2.*f);
    return mix(mix(hash(i), hash(i+vec2(1.,0.)), f.x), mix(hash(i+vec2(0.,1.)), hash(i+vec2(1.,1.)), f.x), f.y); }
  float fbm(vec2 p){ float s = 0., a = .5; for (int i = 0; i < 5; i++){ s += a*noise(p); p = p*2.03 + vec2(1.7, 9.2); a *= .5; } return s; }`;

export function buildScene(scene, camera) {
  const uTime = { value: 0 };
  const uFloor = { value: -0.5 };                                     // 湖底高度，角色載入後對齊腳底
  const sunDir = new THREE.Vector3(-0.45, 0.62, 0.65).normalize();   // 白天的太陽在左上方、觀眾這一側，臉是亮的
  const reflected = [];                                               // 要映在水面上的物件
  const reflect = o => { o.traverse(c => c.layers.enable(REFLECT_LAYER)); reflected.push(o); return o; };

  scene.fog = new THREE.Fog(FOG_COLOR, 60, 320);

  // 天空：上深藍、地平線淡藍；白雲用雜訊畫在高空，順著風慢慢飄
  const sky = reflect(new THREE.Mesh(
    new THREE.SphereGeometry(400, 32, 16),
    new THREE.ShaderMaterial({
      side: THREE.BackSide, depthWrite: false, fog: false,
      uniforms: { uTime, top: { value: SKY_TOP }, horizon: { value: SKY_HORIZON }, sun: { value: SUN_COLOR }, sunDir: { value: sunDir }, wind: { value: WIND } },
      vertexShader: `varying vec3 vDir; void main(){ vDir = normalize(position); gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.); }`,
      fragmentShader: `uniform float uTime; uniform vec3 top, horizon, sun, sunDir; uniform vec2 wind; varying vec3 vDir;
        ${NOISE_GLSL}
        void main(){
          vec3 d = normalize(vDir);
          float h = clamp(d.y, 0., 1.);
          vec3 c = mix(horizon, top, pow(h, .4));
          float s = max(dot(d, sunDir), 0.);
          c += sun * (pow(s, 400.) * 2. + pow(s, 12.) * .15);
          if (d.y > 0.) {
            vec2 uv = d.xz / (d.y + .12) * 1.25 - wind * uTime * .035;
            float n = fbm(uv);
            float cover = smoothstep(.46, .72, n);
            float shade = mix(.8, 1., smoothstep(.45, .85, fbm(uv * 1.8 + 3.1 - wind * uTime * .02)));   // 雲的下緣稍暗
            c = mix(c, vec3(1.) * shade, cover * smoothstep(0., .16, d.y) * .95);
          }
          gl_FragColor = vec4(c, 1.);
          #include <colorspace_fragment>
        }`,
    })));
  scene.add(sky);

  // 預覽：網址加 ?sky=<360 全景圖網址> 時，用全景圖取代程序化的天空與遠山
  const skyUrl = new URLSearchParams(location.search).get('sky');
  const hills = [];
  if (skyUrl) {
    new THREE.TextureLoader().load(skyUrl, tex => {
      tex.mapping = THREE.EquirectangularReflectionMapping;
      tex.colorSpace = THREE.SRGBColorSpace;
      tex.anisotropy = 8;
      sky.material = new THREE.MeshBasicMaterial({ map: tex, side: THREE.BackSide, depthWrite: false, fog: false });
      sky.geometry = new THREE.SphereGeometry(400, 64, 32);
      sky.rotation.y = Number(new URLSearchParams(location.search).get('skyYaw') || 0) * Math.PI / 180;
      sky.position.y = -Number(new URLSearchParams(location.search).get('skyDrop') || 0);
      for (const h of hills) h.visible = false;
    });
  }

  // 遠方長滿青草的緩丘
  const hillMat = new THREE.MeshLambertMaterial({ color: '#7fa65e' });
  for (const [x, z, sx, sy, sz] of [[-70, -150, 110, 14, 40], [40, -170, 120, 18, 45], [120, -140, 70, 10, 35], [-140, -120, 80, 9, 30], [0, -210, 160, 26, 50]]) {
    const hill = new THREE.Mesh(new THREE.SphereGeometry(1, 48, 24), hillMat);
    hill.scale.set(sx, sy, sz); hill.position.set(x, -1, z);
    scene.add(reflect(hill)); hills.push(hill);
  }

  // 湖底：細沙與小石，水面折射下來的光紋（caustics）緩緩流動；靠岸處往上緩升
  const floor = new THREE.Mesh(
    new THREE.PlaneGeometry(90, 80, 1, 160).rotateX(-Math.PI / 2).translate(0, 0, -30),
    new THREE.ShaderMaterial({
      fog: false,
      uniforms: { uTime, uFloor, fogColor: { value: WATER_DEEP } },
      vertexShader: `uniform float uFloor; varying vec3 vWorld;
        void main(){
          vec4 w = modelMatrix * vec4(position, 1.);
          w.y = min(uFloor + (-.04 - uFloor) * smoothstep(${(SHORE_Z - 2.5).toFixed(2)}, ${(SHORE_Z + 0.1).toFixed(2)}, w.z), -.04);
          vWorld = w.xyz;
          gl_Position = projectionMatrix * viewMatrix * w;
        }`,
      fragmentShader: `uniform float uTime; uniform vec3 fogColor; varying vec3 vWorld;
        ${NOISE_GLSL}
        // 光紋：扭曲的正弦網格，取接近零的地方當成亮線
        float caustic(vec2 p, float t){
          vec2 q = p; float c = 0.;
          for (int i = 0; i < 3; i++) {
            q += vec2(sin(q.y * 1.7 + t), cos(q.x * 1.5 - t * .8)) * .35;
            c += abs(sin(q.x * 2.1) * sin(q.y * 2.3));
          }
          return pow(clamp(1. - c / 3., 0., 1.), 6.);
        }
        void main(){
          vec2 p = vWorld.xz;
          float n = fbm(p * 1.6);
          vec3 sand = mix(vec3(.26, .25, .15), vec3(.4, .37, .23), n);
          float pebble = smoothstep(.72, .8, noise(p * 7.));
          sand = mix(sand, vec3(.16, .18, .14), pebble * .7);
          sand += vec3(.9, 1., .85) * caustic(p * 1.1, uTime * .45) * .45;
          sand = mix(sand, fogColor, .25);                                   // 隔著一層水，帶點湖水的青綠
          float d = length(cameraPosition.xz - vWorld.xz);
          vec3 c = mix(sand, fogColor, smoothstep(4., 26., d) * .85);
          gl_FragColor = vec4(c, 1.);
          #include <colorspace_fragment>
        }`,
    }));
  scene.add(floor);

  // 湖面：透明的水，近處看得到湖底與游魚，越遠、越斜看越像鏡子；倒影只映天空、遠山、角色與天鵝
  const water = new Reflector(new THREE.PlaneGeometry(600, 600), {
    textureWidth: 512, textureHeight: 512, clipBias: 0.003, multisample: 0,
    shader: {
      name: 'LakeWater',
      uniforms: { color: { value: null }, tDiffuse: { value: null }, textureMatrix: { value: null }, uTime: { value: 0 }, deep: { value: WATER_DEEP }, fogColor: { value: FOG_COLOR }, wind: { value: WIND } },
      vertexShader: `uniform mat4 textureMatrix; varying vec4 vUv; varying vec3 vWorld;
        void main(){ vUv = textureMatrix * vec4(position, 1.); vec4 w = modelMatrix * vec4(position, 1.); vWorld = w.xyz; gl_Position = projectionMatrix * viewMatrix * w; }`,
      fragmentShader: `uniform sampler2D tDiffuse; uniform float uTime; uniform vec3 deep, fogColor; uniform vec2 wind; varying vec4 vUv; varying vec3 vWorld;
        void main(){
          vec2 p = vWorld.xz;
          float along = dot(p, wind), across = dot(p, vec2(-wind.y, wind.x));
          // 微風吹出的細波，順著風往前推
          float r = sin(along * 1.9 - uTime * 1.3) * .5 + sin(across * 2.4 + along * .6 - uTime * .9) * .3 + sin((along + across) * 4.7 - uTime * 2.2) * .2;
          float r2 = sin(along * 7.3 - uTime * 3.1 + across * 1.7) * .5 + sin(across * 8.1 - uTime * 2.4) * .5;
          vec3 n = normalize(vec3(r * .035 + r2 * .012, 1., r * .028 + r2 * .01));
          vec3 v = normalize(cameraPosition - vWorld);
          float fres = (.03 + .97 * pow(1. - max(dot(n, v), 0.), 5.)) * .65;   // 倒影淡一點，水才透
          vec4 uv = vUv; uv.xy += n.xz * uv.w * .35;
          vec3 refl = texture2DProj(tDiffuse, uv).rgb;
          float d = length(cameraPosition.xz - vWorld.xz);
          float murk = mix(.15, 1., smoothstep(8., 45., d));      // 近處清澈、遠處看不到底
          float a = fres + (1. - fres) * murk;
          vec3 c = (refl * fres + deep * (1. - fres) * murk) / max(a, .001);
          c = mix(c, fogColor, smoothstep(60., 220., d));
          gl_FragColor = vec4(c, a);
          #include <colorspace_fragment>
        }`,
    },
  });
  water.rotateX(-Math.PI / 2);
  water.position.y = WATER_Y;
  water.material.uniforms.uTime = uTime;
  water.material.transparent = true;
  water.material.depthWrite = false;
  scene.add(water);
  const mirrorCam = water.getReflectionCamera(camera);
  mirrorCam.layers.set(REFLECT_LAYER);

  // 角色腳邊的漣漪
  const ripple = ripples(uTime);
  ripple.renderOrder = 1;   // 畫在水面之後，才不會被水蓋住
  scene.add(ripple);

  // 燈光：白天的天光 + 左上方的太陽 + 正面補光
  const sunLight = new THREE.DirectionalLight('#fff4e2', 2.2);
  sunLight.position.copy(sunDir).multiplyScalar(20);
  const fill = new THREE.DirectionalLight('#ffffff', 0.6);
  fill.position.set(1.5, 2.5, 5);
  for (const l of [new THREE.HemisphereLight('#d6e8ff', '#6f8f50', 1.15), sunLight, fill]) scene.add(reflect(l));

  const fish = school();
  scene.add(fish.group);
  const swans = swanPair();
  scene.add(reflect(swans.group));
  const birds = swallows();
  scene.add(birds.group);

  return {
    uTime, sunDir,
    /**
     * 植物要等角色載入後才知道肩膀與腳底高度：綠色蘆葦最高不超過肩膀，從湖底長出；最前景補一條岸邊青草
     * @param {number} maxReedTop 蘆葦（含穗）的最高點，世界座標 y
     * @param {number} floorY 湖底高度（角色腳底），世界座標 y
     */
    addVegetation(maxReedTop, floorY) {
      uFloor.value = floorY;
      fish.setFloor(floorY);
      for (const mesh of reeds(uTime, maxReedTop, floorY)) scene.add(mesh);
      for (const mesh of grass(uTime)) scene.add(mesh);
    },
    /** 漣漪跟著角色移動；走上岸時淡出 */
    setRipple(x, z, strength) {
      ripple.position.x = x; ripple.position.z = z;
      ripple.material.uniforms.uStrength.value = strength;
      ripple.visible = strength > 0.01;
    },
    /** 設定角色撥動植物的點：[{ position: Vector3, radius }]，用 0 號開始的位置（1、2 號留給天鵝） */
    setPushers(list) {
      list.slice(0, 1).forEach((p, i) => PUSH.value[i].set(p.position.x, p.position.y, p.position.z, p.radius));
    },
    /** 每一幀更新游魚、天鵝與燕子 */
    update(t, dt, m) {
      fish.update(t, m);
      swans.update(t, dt * m, PUSH.value);
      birds.update(t, dt * m);
    },
    /** 畫布大小改變時，倒影用一半解析度 */
    resize(w, h) {
      water.getRenderTarget().setSize(Math.max(1, Math.round(w / 2)), Math.max(1, Math.round(h / 2)));
    },
    /** 除錯用：讓天鵝與燕子立刻出現在畫面中間附近 */
    summon() { swans.summon(); birds.summon(); },
    dispose() { water.dispose(); },
  };
}

// 以世界座標的高度計算風吹位移：莖與穗共用同一個公式，穗才會跟著莖一起擺。
// 微風一直吹：一陣陣的風順著風向在草叢上推過去。另外依 PUSH 的點把附近的植物往外撥開（越高處位移越大，根部不動）。
function windShader(mat, uTime, strength, baseY, plantHeight) {
  mat.onBeforeCompile = s => {
    s.uniforms.uTime = uTime;
    s.uniforms.uPush = PUSH;
    s.vertexShader = 'uniform float uTime;\nuniform vec4 uPush[3];\n' + s.vertexShader.replace('#include <project_vertex>', `
      vec4 wp = modelMatrix * instanceMatrix * vec4(transformed, 1.0);
      vec3 base = vec3(instanceMatrix[3]);
      float hh = max(wp.y - (${baseY.toFixed(3)}), 0.);
      float along = dot(base.xz, vec2(${WIND.x.toFixed(3)}, ${WIND.y.toFixed(3)}));
      float gust = sin(uTime*1.1 - along*.45) * .45 + sin(uTime*2.3 + base.x*1.3 + base.z*.7) * .2 + sin(uTime*4.1 + base.z*2.1) * .06 + .55;
      wp.x += gust * ${(strength * WIND.x).toFixed(3)} * hh * hh;
      wp.z += gust * ${(strength * WIND.y).toFixed(3)} * hh * hh;
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

// 綠色蘆葦：一叢一叢從湖底長出，避開角色與鏡頭之間的視線；露出水面的高度介於肩膀的六成到肩膀之間
const PLUME_SCALE = 0.6;                       // 穗長度比例
const PLUME_TOP = 0.46 * PLUME_SCALE;          // 穗頂端高出莖頂的距離（球體 0.26 縮放、往上平移 0.2）
function reeds(uTime, maxTop, floorY) {
  const baseY = floorY;
  let seed = 11; const rnd = () => (seed = (seed * 16807) % 2147483647) / 2147483647;
  const items = [];
  const clumps = 220;
  for (let c = 0; c < clumps; c++) {
    const cx = (rnd() - 0.5) * 36, cz = -18 + rnd() * 20.4;     // 到岸邊草地前為止
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
        tone: 0.85 + rnd() * 0.25,
      });
    }
  }

  const stemGeo = new THREE.PlaneGeometry(0.016, 1, 1, 6).translate(0, 0.5, 0);
  const stemMat = new THREE.MeshLambertMaterial({ color: '#5e8f3c', side: THREE.DoubleSide });
  windShader(stemMat, uTime, 0.06, baseY, maxTop - baseY);
  const stems = new THREE.InstancedMesh(stemGeo, stemMat, items.length);

  // 穗：細長、微微下垂的嫩綠穗
  const plumeGeo = new THREE.SphereGeometry(1, 10, 8).scale(0.026, 0.26, 0.02).translate(0, 0.2, 0);
  const plumeMat = new THREE.MeshLambertMaterial({ color: '#a8cc6e', emissive: '#3c5a1e', emissiveIntensity: 0.3 });
  windShader(plumeMat, uTime, 0.06, baseY, maxTop - baseY);
  const plumes = new THREE.InstancedMesh(plumeGeo, plumeMat, items.length);

  const m = new THREE.Matrix4(), q = new THREE.Quaternion(), s = new THREE.Vector3(1, 1, 1), p = new THREE.Vector3(), e = new THREE.Euler();
  const color = new THREE.Color();
  items.forEach((it, i) => {
    e.set(it.tilt, it.yaw, it.tilt * 0.6);
    q.setFromEuler(e);
    s.set(1, it.h, 1);
    m.compose(p.set(it.x, it.y, it.z), q, s);
    stems.setMatrixAt(i, m);
    stems.setColorAt(i, color.setScalar(it.tone));
    // 穗接在莖的頂端，往一側下垂
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

// 最前景的岸邊青草：只鋪在鏡頭正前方的一條帶狀區域，高度壓低，框住畫面下緣而不擋到人物
function grass(uTime) {
  let seed = 23; const rnd = () => (seed = (seed * 16807) % 2147483647) / 2147483647;
  const group = [];
  const groundY = 0.02;

  const bank = new THREE.Mesh(new THREE.PlaneGeometry(40, 10).rotateX(-Math.PI / 2), new THREE.MeshLambertMaterial({ color: '#6f9a46' }));
  bank.position.set(0, groundY - 0.01, SHORE_Z + 5);   // 岸邊從湖岸線開始
  group.push(bank);

  const count = 9000;
  const items = [];
  for (let i = 0; i < count; i++) {
    const z = SHORE_Z + 0.15 + Math.pow(rnd(), 0.7) * 1.1;
    items.push({ x: -4.5 + rnd() * 10, z, h: 0.26 + rnd() * 0.2, tilt: (rnd() - 0.5) * 0.3, yaw: rnd() * Math.PI, tone: 0.8 + rnd() * 0.35 });
  }

  const bladeGeo = new THREE.PlaneGeometry(0.006, 1, 1, 4).translate(0, 0.5, 0);
  const bladeMat = new THREE.MeshLambertMaterial({ color: '#6aa040', side: THREE.DoubleSide });
  windShader(bladeMat, uTime, 0.2, groundY, 0.55);
  const blades = new THREE.InstancedMesh(bladeGeo, bladeMat, count);

  // 草穗：嫩綠色的小穗
  const tipGeo = new THREE.SphereGeometry(1, 8, 6).scale(0.007, 0.038, 0.007).translate(0, 0.032, 0);
  const tipMat = new THREE.MeshLambertMaterial({ color: '#9ccb5a', emissive: '#2e4a16', emissiveIntensity: 0.25 });
  windShader(tipMat, uTime, 0.2, groundY, 0.55);
  const tips = new THREE.InstancedMesh(tipGeo, tipMat, count);

  const m = new THREE.Matrix4(), q = new THREE.Quaternion(), s = new THREE.Vector3(), p = new THREE.Vector3(), e = new THREE.Euler();
  const color = new THREE.Color();
  items.forEach((it, i) => {
    q.setFromEuler(e.set(it.tilt, it.yaw, it.tilt * 0.5));
    p.set(it.x, groundY, it.z);
    blades.setMatrixAt(i, m.compose(p, q, s.set(1, it.h, 1)));
    blades.setColorAt(i, color.setScalar(it.tone));
    const top = new THREE.Vector3(0, it.h, 0).applyQuaternion(q).add(p);
    q.setFromEuler(e.set(it.tilt + 0.35, it.yaw, it.tilt * 0.5));
    tips.setMatrixAt(i, m.compose(top, q, s.set(1, 1, 1)));
    tips.setColorAt(i, color.setScalar(it.tone));
  });
  blades.frustumCulled = tips.frustumCulled = false;
  group.push(blades, tips);
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
          gl_FragColor = vec4(vec3(1.), a * .35 * uStrength);
        }`,
    }));
  mesh.position.set(0, WATER_Y + 0.004, STAND_Z);
  return mesh;
}

// 水裡的錦鯉：各自繞著一個中心慢慢游，避開角色腳邊；身體左右擺、尾巴拍水
function school() {
  const group = new THREE.Group();
  const bodyGeo = new THREE.SphereGeometry(1, 14, 10).scale(0.032, 0.026, 0.11);
  const tailGeo = new THREE.BufferGeometry().setAttribute('position', new THREE.Float32BufferAttribute([0, 0, 0, 0, 0.045, -0.075, 0, -0.045, -0.075], 3));
  tailGeo.computeVertexNormals();
  const colors = ['#f08a3c', '#f4efe6', '#e9b44c', '#f26b3a', '#f4efe6', '#3e4d55', '#f08a3c', '#e9b44c', '#f4efe6', '#f26b3a'];
  // [中心 x, 中心 z, 半徑, 速度(公尺/秒), 方向]
  const paths = [[-1.7, 0.4, 1.0, 0.32, 1], [-1.4, 0.9, 0.7, 0.26, -1], [1.8, 0.9, 1.2, 0.3, -1], [2.1, 1.5, 0.6, 0.22, 1], [0.6, 1.6, 0.7, 0.28, 1],
    [-2.6, -2.2, 1.1, 0.3, -1], [2.4, -2.6, 1.0, 0.27, 1], [-0.9, 2.0, 0.5, 0.24, -1], [3.0, -0.6, 0.9, 0.3, 1], [-3.2, -0.4, 0.8, 0.26, 1]];
  let floorY = -0.5;
  const fish = paths.map(([cx, cz, r, speed, dir], i) => {
    const color = colors[i % colors.length];
    const mat = new THREE.MeshLambertMaterial({ color, emissive: color, emissiveIntensity: 0.3 });
    const body = new THREE.Mesh(bodyGeo, mat);
    const tail = new THREE.Mesh(tailGeo, new THREE.MeshLambertMaterial({ color, side: THREE.DoubleSide }));
    tail.position.z = -0.1;
    const f = new THREE.Group();
    f.add(body, tail);
    f.scale.setScalar(1.15 + (i % 4) * 0.15);
    group.add(f);
    return { f, tail, cx, cz, r, speed, dir, phase: i * 1.7, depth: 0.35 + (i % 3) * 0.2 };
  });
  return {
    group,
    setFloor(y) { floorY = y; },
    update(t, m) {
      for (const k of fish) {
        const r = k.r * (1 + 0.2 * Math.sin(t * 0.21 + k.phase));
        const a = k.phase + k.dir * t * k.speed * m / k.r;
        const x = k.cx + Math.cos(a) * r, z = k.cz + Math.sin(a) * r;
        // 在水面與湖底之間的某個深度，靠岸變淺時跟著往上
        const bottom = floorY + (-0.04 - floorY) * THREE.MathUtils.smoothstep(z, SHORE_Z - 2.5, SHORE_Z + 0.1);
        k.f.position.set(x, THREE.MathUtils.lerp(bottom + 0.06, WATER_Y - 0.08, k.depth), z);
        const heading = Math.atan2(-Math.sin(a) * k.dir, Math.cos(a) * k.dir);   // 沿切線前進的方向
        k.f.rotation.y = heading + Math.sin(t * 6 + k.phase) * 0.12 * m;
        k.tail.rotation.y = Math.sin(t * 9 + k.phase) * 0.55 * m;
      }
    },
  };
}

// 一對白天鵝：每隔一陣子從畫面一側，橫越角色前方開闊的水面，慢慢游到另一側；游過的地方蘆葦往兩旁讓開
function swanPair() {
  const group = new THREE.Group();
  const make = () => {
    const white = new THREE.MeshLambertMaterial({ color: '#ffffff', emissive: '#e6edf2', emissiveIntensity: 0.25 });
    const s = new THREE.Group();
    const body = new THREE.Mesh(new THREE.SphereGeometry(1, 20, 14).scale(0.16, 0.11, 0.3), white);
    body.position.y = 0.04;
    const wings = new THREE.Mesh(new THREE.SphereGeometry(1, 16, 10).scale(0.15, 0.08, 0.24), white);
    wings.position.set(0, 0.11, -0.05);
    const tail = new THREE.Mesh(new THREE.ConeGeometry(0.06, 0.14, 10).rotateX(-Math.PI / 2 - 0.5), white);
    tail.position.set(0, 0.1, -0.3);
    const neckPath = new THREE.CatmullRomCurve3([new THREE.Vector3(0, 0.1, 0.2), new THREE.Vector3(0, 0.28, 0.27), new THREE.Vector3(0, 0.44, 0.2), new THREE.Vector3(0, 0.53, 0.25)]);
    const neck = new THREE.Mesh(new THREE.TubeGeometry(neckPath, 20, 0.028, 10), white);
    const head = new THREE.Group();
    head.position.set(0, 0.54, 0.28);
    head.add(new THREE.Mesh(new THREE.SphereGeometry(1, 14, 10).scale(0.042, 0.04, 0.065), white));
    const beak = new THREE.Mesh(new THREE.ConeGeometry(0.016, 0.07, 10).rotateX(Math.PI / 2 + 0.35), new THREE.MeshLambertMaterial({ color: '#e8772e' }));
    beak.position.set(0, -0.012, 0.085);
    const knob = new THREE.Mesh(new THREE.SphereGeometry(0.014, 8, 6), new THREE.MeshLambertMaterial({ color: '#1f1f1f' }));
    knob.position.set(0, 0.004, 0.055);
    head.add(beak, knob);
    s.add(body, wings, tail, neck, head);
    s.scale.setScalar(1.2);
    s.visible = false;
    group.add(s);
    return { s, head };
  };
  const pair = [make(), make()];
  const SPEED = 0.3, HALF = 11;
  let wait = 5, x = 0, z = 1.5, dir = 1, active = false;
  const start = () => { active = true; dir = Math.random() < 0.5 ? 1 : -1; x = -dir * HALF; z = 0.1 + Math.random() * 0.9; };
  return {
    group,
    summon() { start(); x = -dir * 2.5; },
    update(t, dt, push) {
      if (!active) {
        wait -= dt;
        if (wait <= 0) start();
        for (const p of pair) p.s.visible = false;
        push[1].w = push[2].w = 0;
        if (!active) return;
      }
      x += dir * SPEED * dt;
      pair.forEach((p, i) => {
        const px = x - dir * i * 1.5, pz = z - i * 0.45;
        p.s.visible = true;
        p.s.position.set(px, WATER_Y + Math.sin(t * 1.3 + i) * 0.008, pz);
        p.s.rotation.y = dir > 0 ? Math.PI / 2 : -Math.PI / 2;
        p.s.rotation.z = Math.sin(t * 0.9 + i * 2) * 0.03;
        p.head.rotation.x = Math.sin(t * 0.7 + i) * 0.12;
        push[1 + i].set(px, 0, pz, 0.7);
      });
      if (Math.abs(x) > HALF + 2) { active = false; wait = 25 + Math.random() * 25; }
    },
  };
}

// 燕子：每隔一陣子一小群從天空掠過，邊拍翅邊滑翔，各自在群裡上下穿梭
function swallows() {
  const group = new THREE.Group();
  const dark = new THREE.MeshLambertMaterial({ color: '#1e2a3a', side: THREE.DoubleSide });
  const wingGeo = new THREE.BufferGeometry().setAttribute('position', new THREE.Float32BufferAttribute([0, 0, 0.03, 0.3, 0, -0.06, 0, 0, -0.04, 0.3, 0, -0.06, 0.12, 0, -0.07, 0, 0, -0.04], 3));
  wingGeo.computeVertexNormals();
  const tailGeo = new THREE.BufferGeometry().setAttribute('position', new THREE.Float32BufferAttribute([0, 0, -0.07, 0.06, 0, -0.2, 0.02, 0, -0.08, 0, 0, -0.07, -0.02, 0, -0.08, -0.06, 0, -0.2], 3));
  tailGeo.computeVertexNormals();
  const bodyGeo = new THREE.SphereGeometry(1, 10, 8).scale(0.035, 0.03, 0.1);
  const N = 9;
  const birds = Array.from({ length: N }, (_, i) => {
    const b = new THREE.Group();
    const left = new THREE.Mesh(wingGeo, dark);
    const right = new THREE.Mesh(wingGeo, dark);
    right.scale.x = -1;
    b.add(new THREE.Mesh(bodyGeo, dark), left, right, new THREE.Mesh(tailGeo, dark));
    b.scale.setScalar(1.6);
    group.add(b);
    return { b, left, right, off: new THREE.Vector3((Math.random() - 0.5) * 6, (Math.random() - 0.5) * 2.2, (Math.random() - 0.5) * 5), phase: i * 0.9 + Math.random() };
  });
  group.visible = false;
  const SPEED = 5.5, HALF = 30;
  let wait = 9, x = 0, y = 7, z = -18, dir = 1;
  const start = () => {
    group.visible = true;
    dir = Math.random() < 0.5 ? 1 : -1;
    x = -dir * HALF; y = 4.2 + Math.random() * 2.2; z = -12 - Math.random() * 6;   // 遠山上方的那一條天空
  };
  const p0 = new THREE.Vector3(), p1 = new THREE.Vector3();
  const at = (k, cx, tt, out) => out.set(
    cx + k.off.x + Math.sin(tt * 1.3 + k.phase) * 0.8,
    y + k.off.y + Math.sin(tt * 1.7 + k.phase) * 0.6,
    z + k.off.z + Math.cos(tt * 1.1 + k.phase) * 0.9);
  return {
    group,
    summon() { start(); x = -dir * 4; },
    update(t, dt) {
      if (!group.visible) {
        wait -= dt;
        if (wait > 0) return;
        start();
      }
      x += dir * SPEED * dt;
      for (const k of birds) {
        at(k, x, t, p0);
        at(k, x + dir * SPEED * 0.05, t + 0.05, p1);
        k.b.position.copy(p0);
        k.b.lookAt(p1);
        // 拍幾下翅膀，再展翅滑翔一會兒
        const flap = Math.sin(t * 15 + k.phase) * 0.7 * THREE.MathUtils.smoothstep(Math.sin(t * 1.4 + k.phase), -0.3, 0.3);
        k.left.rotation.z = flap; k.right.rotation.z = -flap;
      }
      if (Math.abs(x) > HALF + 8) { group.visible = false; wait = 22 + Math.random() * 30; }
    },
  };
}
