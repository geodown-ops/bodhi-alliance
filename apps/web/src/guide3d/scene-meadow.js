// 晨霧草原：Skybox AI 生成的 360° 全景圖（Geodown 付費帳號產出），Sunny 捧著熱茶站在圖中的河裡。
// 全景圖下半部投影成平地（GroundedSkybox）；遠景鏡頭放在全景圖的拍攝點，看到的畫面就和 Skybox 網站上一樣，
// Sunny 站在河流分岔處偏左的水裡。其他物件（湖、蘆葦、魚、鳥）都不放。
// 河道用一張遮罩圖標出，著色器只在河道上讓倒影輕輕晃動、亮紋順著水流漂過，做出流動感。
// 提供與 scene.js 相同的介面，avatar.js 依時間挑一個。
import * as THREE from 'three';
import { GroundedSkybox } from 'three/addons/objects/GroundedSkybox.js';
import { plantBodhiTree } from './bodhi-tree.js';

const PANORAMA = '/scenes/meadow-dawn.jpg';          // 4K：先載入，很快就有畫面
const PANORAMA_8K = '/scenes/meadow-dawn-8k.jpg';    // 8K 原圖：顯示卡撐得住時接著換上，景色更清晰
const RIVER_MASK = '/scenes/meadow-dawn-river.png';   // 白色 = 河道（1024×512，與全景圖對齊）
const EYE_HEIGHT = 1.2;      // 鏡頭（全景圖拍攝點）離地高度：越低，Sunny 離鏡頭越近、看起來越大
// 遠景構圖：比照 Geodown 在 Skybox 網站上截的角度（山坡上的樹、河流在草地前分岔）
const VIEW_U = 0.078;        // 畫面中心對著全景圖的哪一個橫向位置（0～1）
const VIEW_PITCH = 2;        // 鏡頭仰角（度）：頁面把畫面往上挪了一些（setViewOffset），這裡稍微抬頭補回來
// Sunny 站的位置：河流分岔處偏左的水裡（全景圖座標 u, v）
const RIVER_SPOT = [0.043, 0.549];
const CAM_Z = 6.4;
// 菩提樹：種在畫面中間、Sunny 右前方河對岸的草地上（全景圖座標 u, v）與地面以上的高度
const TREE_SPOT = [0.085, 0.528];
const TREE_HEIGHT = 5;           // 遠景鏡頭的位置（和湖景相同）

const VNOISE = `
  float vhash(vec2 p) { return fract(sin(dot(p, vec2(127.1, 311.7))) * 43758.5453); }
  float vnoise(vec2 p) {
    vec2 i = floor(p), f = fract(p); f = f * f * (3. - 2. * f);
    return mix(mix(vhash(i), vhash(i + vec2(1, 0)), f.x), mix(vhash(i + vec2(0, 1)), vhash(i + vec2(1, 1)), f.x), f.y);
  }`;

/** 全景圖上一點投影到地面後，相對拍攝點的水平位置（GroundedSkybox 的球面 UV 對應，未旋轉） */
function groundPoint(u, v, height) {
  const phi = u * Math.PI * 2, below = (v - 0.5) * Math.PI;
  const dist = height / Math.tan(below);
  return new THREE.Vector2(-Math.cos(phi) * dist, -Math.sin(phi) * dist);
}

export function buildMeadowScene(scene, camera, renderer) {
  const uTime = { value: 0 };
  const sunDir = new THREE.Vector3(-0.3, 0.5, 0.8).normalize();   // 晨光從觀眾這一側照到臉上
  const q = new URLSearchParams(location.search);
  const height = Number(q.get('eye')) || EYE_HEIGHT;
  const spotUV = [Number(q.get('spotU')) || RIVER_SPOT[0], Number(q.get('spotV')) || RIVER_SPOT[1]];
  // 轉動全景圖，讓構圖中心正對 -z（鏡頭看的方向）
  const phi = VIEW_U * Math.PI * 2;
  const yaw = Math.PI - Math.atan2(-Math.cos(phi), -Math.sin(phi));
  const spot = groundPoint(spotUV[0], spotUV[1], height).rotateAround(new THREE.Vector2(), -yaw);
  const standAt = { x: spot.x, z: CAM_Z + spot.y };
  const pitch = THREE.MathUtils.degToRad(VIEW_PITCH);
  const wide = {
    pos: new THREE.Vector3(0, height, CAM_Z),
    look: new THREE.Vector3(0, height + Math.sin(pitch) * 10, CAM_Z - Math.cos(pitch) * 10),
  };

  scene.background = new THREE.Color('#d9a6b8');   // 全景圖載入前的粉色晨空
  let sky = null;
  const loader = new THREE.TextureLoader();
  const mask = loader.load(RIVER_MASK);
  const anisotropy = renderer?.capabilities.getMaxAnisotropy() ?? 8;
  loader.load(PANORAMA, tex => {
    tex.colorSpace = THREE.SRGBColorSpace;
    tex.anisotropy = anisotropy;
    sky = new GroundedSkybox(tex, height, 300);
    // 地面會寫入深度：站在河裡時，水面以下的小腿被河面蓋住
    sky.material.depthWrite = true;
    sky.renderOrder = -1;
    sky.material.onBeforeCompile = shader => {
      shader.uniforms.uTime = uTime;
      shader.uniforms.uRiver = { value: mask };
      shader.fragmentShader = shader.fragmentShader
        .replace('#include <map_pars_fragment>', '#include <map_pars_fragment>\nuniform float uTime; uniform sampler2D uRiver;\n' + VNOISE)
        .replace('#include <map_fragment>', `
          vec2 uv = vMapUv;
          float river = texture2D(uRiver, uv).r;
          if (river > 0.01) {
            // 水流沿著河道（全景圖的橫向）往前推：倒影隨著不規則的細波輕輕晃動
            vec2 p = vec2(uv.x * 1400., uv.y * 2600.);
            float n1 = vnoise(p + vec2(-uTime * 1.1, 0.)), n2 = vnoise(p * 2.3 + vec2(-uTime * 1.9, 3.7));
            uv += river * vec2((n2 - .5) * .00012, (n1 - .5) * .00045);
          }
          vec4 sampledDiffuseColor = texture2D(map, uv);
          if (river > 0.01) {
            // 亮點：細碎的反光順著水流閃動漂過
            float g = vnoise(vec2(uv.x * 1100. - uTime * 1.6, uv.y * 2400.)) * vnoise(vec2(uv.x * 420. - uTime * .8, uv.y * 900. + 5.));
            sampledDiffuseColor.rgb += river * smoothstep(.5, .8, g) * .1 * vec3(1., .96, .98);
          }
          diffuseColor *= sampledDiffuseColor;`);
    };
    // 全景圖的中心就是遠景鏡頭的位置
    sky.position.set(0, height - 0.01, CAM_Z);
    sky.rotation.y = yaw;
    scene.add(sky);
    if ((renderer?.capabilities.maxTextureSize ?? 0) >= 8192) {
      loader.load(PANORAMA_8K, hd => {
        if (!sky) return;
        hd.colorSpace = THREE.SRGBColorSpace;
        hd.anisotropy = anisotropy;
        const old = sky.material.map;
        sky.material.map = hd;
        old.dispose();
      });
    }
  });

  const treeUV = [Number(q.get('treeU')) || TREE_SPOT[0], Number(q.get('treeV')) || TREE_SPOT[1]];
  const treeAt = groundPoint(treeUV[0], treeUV[1], height).rotateAround(new THREE.Vector2(), -yaw);
  const tree = plantBodhiTree(scene, {
    x: treeAt.x, z: CAM_Z + treeAt.y, height: Number(q.get('treeH')) || TREE_HEIGHT, yaw: 0.6, uTime,
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
  ripple.position.set(standAt.x, 0.005, standAt.z);
  scene.add(ripple);

  // 燈光：粉橘色的晨光 + 天空的柔光 + 正面補光
  const sunLight = new THREE.DirectionalLight('#ffe2d0', 1.9);
  sunLight.position.copy(sunDir).multiplyScalar(20);
  const fill = new THREE.DirectionalLight('#ffffff', 0.6);
  fill.position.set(1.5, 2.5, 5);
  for (const l of [new THREE.HemisphereLight('#f3d6e4', '#7d9458', 1.2), sunLight, fill]) scene.add(l);

  return {
    uTime, sunDir,
    standAt, wide,      // Sunny 的位置、遠景鏡頭
    wideHFov: 66,       // 遠景的水平視角（度），和 Skybox 網站上的畫面一樣寬
    onGround: true,     // 站在投影的地面上（河裡），不是程式畫的湖
    avatarScale: 0.5,   // 依全景圖的比例，人物縮成一半
    wadeDepth: 0.55,    // 河水淹到小腿約一半（膝蓋高度的比例）
    addVegetation() {},
    setRipple(x, z) { ripple.position.x = x; ripple.position.z = z; },
    setPushers() {},
    update(t, dt) { tree.update(t, dt); },
    resize() {},
    summon() {},
    dispose() { sky?.material.map?.dispose(); mask.dispose(); tree.dispose(); },
  };
}
