// 晨霧草原：Skybox AI 生成的 360° 全景圖（Geodown 付費帳號產出），只有 Sunny 捧著熱茶站在草地上。
// 全景圖下半部投影成平地（GroundedSkybox），人物才像真的站在草原上；其他物件（湖、蘆葦、魚、鳥）都不放。
// 提供與 scene.js 相同的介面，avatar.js 依時間挑一個。
import * as THREE from 'three';
import { GroundedSkybox } from 'three/addons/objects/GroundedSkybox.js';

const PANORAMA = '/scenes/meadow-dawn.jpg';
const CAPTURE_HEIGHT = 4;    // 全景圖拍攝點離地高度：越大地面越寬闊
const YAW = 0;               // 全景圖轉向（度），決定 Sunny 身後是哪一片風景

export function buildMeadowScene(scene) {
  const uTime = { value: 0 };
  const sunDir = new THREE.Vector3(-0.3, 0.5, 0.8).normalize();   // 晨光從觀眾這一側照到臉上
  const q = new URLSearchParams(location.search);
  const height = Number(q.get('skyHeight')) || CAPTURE_HEIGHT;
  const yaw = (q.has('skyYaw') ? Number(q.get('skyYaw')) : YAW) * Math.PI / 180;

  scene.background = new THREE.Color('#d9a6b8');   // 全景圖載入前的粉色晨空
  let sky = null;
  new THREE.TextureLoader().load(PANORAMA, tex => {
    tex.mapping = THREE.EquirectangularReflectionMapping;
    tex.colorSpace = THREE.SRGBColorSpace;
    tex.anisotropy = 8;
    sky = new GroundedSkybox(tex, height, 300);
    sky.position.y = height - 0.01;
    sky.rotation.y = yaw;
    scene.add(sky);
  });

  // 腳下柔和的影子，讓人物落地
  const shadowTex = new THREE.CanvasTexture((() => {
    const c = document.createElement('canvas'); c.width = c.height = 128;
    const g = c.getContext('2d'), r = g.createRadialGradient(64, 64, 0, 64, 64, 64);
    r.addColorStop(0, 'rgba(30,40,20,0.55)'); r.addColorStop(1, 'rgba(30,40,20,0)');
    g.fillStyle = r; g.fillRect(0, 0, 128, 128); return c;
  })());
  const shadow = new THREE.Mesh(new THREE.PlaneGeometry(1.1, 0.8), new THREE.MeshBasicMaterial({ map: shadowTex, transparent: true, depthWrite: false }));
  shadow.rotation.x = -Math.PI / 2;
  shadow.position.y = 0.005;
  shadow.renderOrder = 1;
  scene.add(shadow);

  // 燈光：粉橘色的晨光 + 天空的柔光 + 正面補光
  const sunLight = new THREE.DirectionalLight('#ffe2d0', 1.9);
  sunLight.position.copy(sunDir).multiplyScalar(20);
  const fill = new THREE.DirectionalLight('#ffffff', 0.6);
  fill.position.set(1.5, 2.5, 5);
  for (const l of [new THREE.HemisphereLight('#f3d6e4', '#7d9458', 1.2), sunLight, fill]) scene.add(l);

  return {
    uTime, sunDir,
    onGround: true,   // 站在草地上，不是站在水裡
    addVegetation() {},
    setRipple(x, z) { shadow.position.x = x; shadow.position.z = z; },
    setPushers() {},
    update() {},
    resize() {},
    summon() {},
    dispose() { sky?.material.map?.dispose(); },
  };
}
