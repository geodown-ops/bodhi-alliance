// 晨霧草原：Skybox AI 生成的 360° 全景圖（Geodown 付費帳號產出），Sunny 捧著熱茶站在圖中的河裡。
// 全景圖下半部投影成平地（GroundedSkybox），再把投影移到河道正好落在 Sunny 腳下；其他物件（湖、蘆葦、魚、鳥）都不放。
// 河道用一張遮罩圖標出，著色器只在河道上讓倒影輕輕晃動、亮紋順著水流漂過，做出流動感。
// 提供與 scene.js 相同的介面，avatar.js 依時間挑一個。
import * as THREE from 'three';
import { GroundedSkybox } from 'three/addons/objects/GroundedSkybox.js';
import { STAND_Z } from './scene.js';

const PANORAMA = '/scenes/meadow-dawn.jpg';
const RIVER_MASK = '/scenes/meadow-dawn-river.png';   // 白色 = 河道（1024×512，與全景圖對齊）
const CAPTURE_HEIGHT = 4;    // 全景圖拍攝點離地高度：越大地面越寬闊
const YAW = 0;               // 全景圖轉向（度），決定 Sunny 身後是哪一片風景
// 河道上要讓 Sunny 站的那一點（全景圖座標 u, v，0～1）：近處那段河、旭日在右後方
const RIVER_SPOT = [0.234, 0.645];
const SPOT_NUDGE = -1.4;       // 投影的平地與原圖略有偏差，往前微調（公尺）讓腳正好落在河中央

/** 全景圖上一點投影到地面後，相對拍攝點的水平位置（GroundedSkybox 的球面 UV 對應，未旋轉） */
function groundPoint(u, v, height) {
  const phi = u * Math.PI * 2, below = (v - 0.5) * Math.PI;
  const dist = height / Math.tan(below);
  return new THREE.Vector2(-Math.cos(phi) * dist, -Math.sin(phi) * dist);
}

export function buildMeadowScene(scene) {
  const uTime = { value: 0 };
  const sunDir = new THREE.Vector3(-0.3, 0.5, 0.8).normalize();   // 晨光從觀眾這一側照到臉上
  const q = new URLSearchParams(location.search);
  const height = Number(q.get('skyHeight')) || CAPTURE_HEIGHT;
  const yaw = (q.has('skyYaw') ? Number(q.get('skyYaw')) : YAW) * Math.PI / 180;

  scene.background = new THREE.Color('#d9a6b8');   // 全景圖載入前的粉色晨空
  let sky = null;
  const loader = new THREE.TextureLoader();
  const mask = loader.load(RIVER_MASK);
  loader.load(PANORAMA, tex => {
    tex.colorSpace = THREE.SRGBColorSpace;
    tex.anisotropy = 8;
    sky = new GroundedSkybox(tex, height, 300);
    // 地面會寫入深度：站在河裡時，水面以下的小腿被河面蓋住
    sky.material.depthWrite = true;
    sky.renderOrder = -1;
    sky.material.onBeforeCompile = shader => {
      shader.uniforms.uTime = uTime;
      shader.uniforms.uRiver = { value: mask };
      shader.fragmentShader = shader.fragmentShader
        .replace('#include <map_pars_fragment>', '#include <map_pars_fragment>\nuniform float uTime; uniform sampler2D uRiver;')
        .replace('#include <map_fragment>', `
          vec2 uv = vMapUv;
          float river = texture2D(uRiver, uv).r;
          if (river > 0.01) {
            // 水流沿著河道（全景圖的橫向）往前推：倒影細細晃動
            float w1 = sin(uv.x * 900. - uTime * 2.2 + sin(uv.y * 700.) * 1.5);
            float w2 = sin(uv.x * 2300. + uv.y * 1500. - uTime * 3.7);
            uv += river * vec2(w2 * .00008, (w1 * .7 + w2 * .3) * .00022);
          }
          vec4 sampledDiffuseColor = texture2D(map, uv);
          if (river > 0.01) {
            // 亮紋：一道道淡淡的反光順著水流漂過
            float s = sin(uv.x * 520. - uTime * 1.6 + sin(uv.y * 900. + uTime * .7) * 2.2);
            float glint = pow(max(s, 0.), 16.) * .14 + pow(max(sin(uv.x * 1300. - uTime * 2.4 + sin(uv.y * 600.) * 3.), 0.), 40.) * .08;
            sampledDiffuseColor.rgb += river * glint * vec3(1., .95, .97);
          }
          diffuseColor *= sampledDiffuseColor;`);
    };
    // 把投影移過來，讓選定的那一點河道落在 Sunny 腳下
    const spot = groundPoint(RIVER_SPOT[0], RIVER_SPOT[1], height).rotateAround(new THREE.Vector2(), -yaw);
    sky.position.set(-spot.x + Number(q.get('skyX') || 0), height - 0.01, STAND_Z - spot.y + SPOT_NUDGE + Number(q.get('skyZ') || 0));
    sky.rotation.y = yaw;
    scene.add(sky);
  });

  // 腳邊的漣漪：一圈圈往外擴散、慢慢淡出
  const ripple = new THREE.Mesh(
    new THREE.PlaneGeometry(1.4, 1.4),
    new THREE.ShaderMaterial({
      transparent: true, depthWrite: false,
      uniforms: { uTime },
      vertexShader: 'varying vec2 vUv; void main(){ vUv = uv; gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.); }',
      fragmentShader: `uniform float uTime; varying vec2 vUv;
        void main(){
          float r = length(vUv - .5) * 2.;
          float a = 0.;
          for (int i = 0; i < 3; i++) {
            float p = fract(uTime * .22 + float(i) / 3.);
            a += smoothstep(.035, 0., abs(r - (.12 + p * .85))) * (1. - p) * .5;
          }
          a += smoothstep(.16, .08, r) * .25;   // 小腿入水處的一圈白色水花
          gl_FragColor = vec4(vec3(1., .96, .97), a * smoothstep(1., .8, r));
        }`,
    }));
  ripple.rotation.x = -Math.PI / 2;
  ripple.position.set(0, 0.005, STAND_Z);
  scene.add(ripple);

  // 燈光：粉橘色的晨光 + 天空的柔光 + 正面補光
  const sunLight = new THREE.DirectionalLight('#ffe2d0', 1.9);
  sunLight.position.copy(sunDir).multiplyScalar(20);
  const fill = new THREE.DirectionalLight('#ffffff', 0.6);
  fill.position.set(1.5, 2.5, 5);
  for (const l of [new THREE.HemisphereLight('#f3d6e4', '#7d9458', 1.2), sunLight, fill]) scene.add(l);

  return {
    uTime, sunDir,
    onGround: true,     // 站在投影的地面上（河裡），不是程式畫的湖
    avatarScale: 0.5,   // 依全景圖的比例，人物縮成一半
    wadeDepth: 0.55,    // 河水淹到小腿約一半（膝蓋高度的比例）
    addVegetation() {},
    setRipple(x, z) { ripple.position.x = x; ripple.position.z = z; },
    setPushers() {},
    update() {},
    resize() {},
    summon() {},
    dispose() { sky?.material.map?.dispose(); mask.dispose(); },
  };
}
