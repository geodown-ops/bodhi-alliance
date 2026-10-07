// 全景圖場景的共用骨架：Skybox AI 生成的 360° 全景圖（Geodown 付費帳號產出），Sunny 捧著熱茶站在圖中的河裡。
// 全景圖下半部投影成平地（GroundedSkybox）；遠景鏡頭放在全景圖的拍攝點，看到的畫面就和 Skybox 網站上一樣。
// 河道用一張遮罩圖標出，著色器只在河道上讓倒影輕輕晃動、亮紋順著水流漂過，做出流動感。
// 晨霧草原（scene-meadow.js）、松林雪山（scene-pines.js）各自填入圖檔、構圖、燈光，提供與 scene.js 相同的介面。
import * as THREE from 'three';
import { GroundedSkybox } from 'three/addons/objects/GroundedSkybox.js';

const CAM_Z = 6.4;           // 遠景鏡頭的位置（和湖景相同）

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

// 顏色：字串是一般的 CSS 色碼；[r, g, b] 直接當著色器裡的數值用
const css = c => Array.isArray(c) ? new THREE.Color().setRGB(...c) : new THREE.Color(c);

/**
 * cfg：
 *   panorama / panorama8k / riverMask  圖檔路徑（4K 先載入，8K 在顯示卡撐得住時換上；遮罩白色 = 河道）
 *   eyeHeight   鏡頭（拍攝點）離地高度：越低，Sunny 離鏡頭越近、看起來越大
 *   viewU, viewPitch  畫面中心對著全景圖的橫向位置（0～1）與仰角（度）
 *   spot        Sunny 站的位置（全景圖座標 [u, v]）
 *   wideHFov    遠景的水平視角（度）
 *   avatarScale 人物的大小（依全景圖的比例，預設 0.5）
 *   background  全景圖載入前的天空顏色
 *   glint       水面反光的顏色（[r, g, b]）
 *   sunDir, sun: [顏色, 強度], hemi: [天色, 地色, 強度], fill: 正面補光強度
 *   ripple      腳邊漣漪的顏色（[r, g, b]）
 *   extras(ctx) 額外的物件（例如菩提樹、鹿）：ctx = { scene, uTime, height, at(u, v), riverAt(x, z), riverReady, standAt, camZ, q }，
 *               回傳 { update(t, dt), dispose() }
 */
export function buildPanoramaScene(scene, camera, renderer, cfg) {
  const uTime = { value: 0 };
  const sunDir = cfg.sunDir.clone().normalize();
  const q = new URLSearchParams(location.search);
  const height = Number(q.get('eye')) || cfg.eyeHeight;
  const spotUV = [Number(q.get('spotU')) || cfg.spot[0], Number(q.get('spotV')) || cfg.spot[1]];
  // 轉動全景圖，讓構圖中心正對 -z（鏡頭看的方向）
  const phi = (Number(q.get('viewU')) || cfg.viewU) * Math.PI * 2;
  const yaw = Math.PI - Math.atan2(-Math.cos(phi), -Math.sin(phi));
  /** 全景圖座標 (u, v) → 場景裡的地面位置 */
  const at = (u, v) => {
    const p = groundPoint(u, v, height).rotateAround(new THREE.Vector2(), -yaw);
    return { x: p.x, z: CAM_Z + p.y };
  };
  const standAt = at(spotUV[0], spotUV[1]);
  const avatarScale = Number(q.get('size')) || cfg.avatarScale || 0.5;
  const pitch = THREE.MathUtils.degToRad(q.has('pitch') ? Number(q.get('pitch')) : cfg.viewPitch);
  const wide = {
    pos: new THREE.Vector3(0, height, CAM_Z),
    look: new THREE.Vector3(0, height + Math.sin(pitch) * 10, CAM_Z - Math.cos(pitch) * 10),
  };

  scene.background = css(cfg.background);
  let sky = null;
  const loader = new THREE.TextureLoader();
  // 河道遮罩同時讀成像素資料，讓鹿之類會走動的東西知道哪裡是河
  let riverData = null, riverLoaded;
  const riverReady = new Promise(r => { riverLoaded = r; });
  const mask = loader.load(cfg.riverMask, tex => {
    const img = tex.image, c = document.createElement('canvas');
    c.width = img.width; c.height = img.height;
    const g = c.getContext('2d', { willReadFrequently: true });
    g.drawImage(img, 0, 0);
    riverData = { w: c.width, h: c.height, px: g.getImageData(0, 0, c.width, c.height).data };
    riverLoaded();
  });
  /** 場景裡的地面位置 → 河道遮罩的值（0 = 陸地，1 = 河）；遮罩還沒載入時一律當成河 */
  const riverAt = (x, z) => {
    if (!riverData) return 1;
    const p = new THREE.Vector2(x, z - CAM_Z).rotateAround(new THREE.Vector2(), yaw);
    const u = ((Math.atan2(-p.y, -p.x) / (Math.PI * 2)) % 1 + 1) % 1;
    const v = 0.5 + Math.atan(height / Math.max(p.length(), 1e-3)) / Math.PI;
    const { w, h, px } = riverData;
    return px[(Math.min(h - 1, Math.floor(v * h)) * w + Math.min(w - 1, Math.floor(u * w))) * 4] / 255;
  };
  const anisotropy = renderer?.capabilities.getMaxAnisotropy() ?? 8;
  const glint = css(cfg.glint);
  loader.load(cfg.panorama, tex => {
    tex.colorSpace = THREE.SRGBColorSpace;
    tex.anisotropy = anisotropy;
    sky = new GroundedSkybox(tex, height, 300, 256);   // 球面切細一點：地平線附近的近物（例如大樹幹）不會出現鋸齒接縫
    // 地面會寫入深度：站在河裡時，水面以下的小腿被河面蓋住
    sky.material.depthWrite = true;
    sky.renderOrder = -1;
    sky.material.onBeforeCompile = shader => {
      shader.uniforms.uTime = uTime;
      shader.uniforms.uRiver = { value: mask };
      shader.uniforms.uGlint = { value: glint };
      shader.fragmentShader = shader.fragmentShader
        .replace('#include <map_pars_fragment>', '#include <map_pars_fragment>\nuniform float uTime; uniform sampler2D uRiver; uniform vec3 uGlint;\n' + VNOISE)
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
            sampledDiffuseColor.rgb += river * smoothstep(.5, .8, g) * .1 * uGlint;
          }
          diffuseColor *= sampledDiffuseColor;`);
    };
    // 全景圖的中心就是遠景鏡頭的位置
    sky.position.set(0, height - 0.01, CAM_Z);
    sky.rotation.y = yaw;
    scene.add(sky);
    if ((renderer?.capabilities.maxTextureSize ?? 0) >= 8192) {
      loader.load(cfg.panorama8k, hd => {
        if (!sky) return;
        hd.colorSpace = THREE.SRGBColorSpace;
        hd.anisotropy = anisotropy;
        const old = sky.material.map;
        sky.material.map = hd;
        old.dispose();
      });
    }
  });

  const extras = cfg.extras?.({ scene, uTime, height, at, riverAt, riverReady, standAt, camZ: CAM_Z, q });

  // 腳邊的漣漪：一圈圈往外擴散、慢慢淡出
  const ripple = new THREE.Mesh(
    new THREE.PlaneGeometry(1.4, 1.4),
    new THREE.ShaderMaterial({
      transparent: true, depthWrite: false,
      uniforms: { uTime, uColor: { value: css(cfg.ripple) } },
      vertexShader: 'varying vec2 vUv; void main(){ vUv = uv; gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.); }',
      fragmentShader: `uniform float uTime; uniform vec3 uColor; varying vec2 vUv;
        void main(){
          float r = length(vUv - .5) * 2.;
          float a = 0.;
          for (int i = 0; i < 3; i++) {
            float p = fract(uTime * .22 + float(i) / 3.);
            a += smoothstep(.035, 0., abs(r - (.12 + p * .85))) * (1. - p) * .5;
          }
          a += smoothstep(.16, .08, r) * .25;   // 小腿入水處的一圈白色水花
          gl_FragColor = vec4(uColor, a * smoothstep(1., .8, r));
        }`,
    }));
  ripple.rotation.x = -Math.PI / 2;
  ripple.scale.setScalar(avatarScale / 0.5);   // 漣漪跟著人物大小
  ripple.position.set(standAt.x, 0.005, standAt.z);
  scene.add(ripple);

  // 燈光：主光（與全景圖的光線一致）+ 天空的柔光 + 正面補光
  const sunLight = new THREE.DirectionalLight(cfg.sun[0], cfg.sun[1]);
  sunLight.position.copy(sunDir).multiplyScalar(20);
  const fill = new THREE.DirectionalLight('#ffffff', cfg.fill);
  fill.position.set(1.5, 2.5, 5);
  for (const l of [new THREE.HemisphereLight(...cfg.hemi), sunLight, fill]) scene.add(l);

  return {
    uTime, sunDir,
    standAt, wide,      // Sunny 的位置、遠景鏡頭
    wideHFov: cfg.wideHFov,
    onGround: true,     // 站在投影的地面上（河裡），不是程式畫的湖
    avatarScale,
    wadeDepth: 0.55,    // 河水淹到小腿約一半（膝蓋高度的比例）
    addVegetation() {},
    setRipple(x, z) { ripple.position.x = x; ripple.position.z = z; },
    setPushers() {},
    update(t, dt) { extras?.update(t, dt); },
    resize() {},
    summon() {},
    dispose() { sky?.material.map?.dispose(); mask.dispose(); extras?.dispose(); },
  };
}
