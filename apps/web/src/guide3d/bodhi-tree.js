// 菩提樹：Geodown 提供的 SpeedTree 模型（樹皮減面後轉成 meshopt 壓縮的 GLB），
// 樹枝隨風輕擺、葉子細碎顫動，樹冠下有菩提葉一片片飄落。
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { MeshoptDecoder } from 'three/addons/libs/meshopt_decoder.module.js';

const BASE = '/models/bodhi-tree/';
const MODEL_HEIGHT = 36;     // 模型地面以上的高度（模型單位）
const FALLING = 70;          // 同時飄落的葉子數
const LIFT = 0.55;           // 樹枝往上抬：離樹幹每 1 單位，抬高 0.55 單位（原模型的枝條較平展、下垂）
const SPREAD = 0.3;          // 樹冠同時向四周展開 30%

/** 把（量化過的）模型轉回一般座標，並讓樹枝往上伸展；回傳抬高後的樹高 */
function raiseBranches(mesh) {
  const g = mesh.geometry, src = g.attributes.position, n = src.count;
  const pos = new Float32Array(n * 3), v = new THREE.Vector3();
  mesh.updateMatrix();
  let top = 0;
  for (let i = 0; i < n; i++) {
    v.fromBufferAttribute(src, i).applyMatrix4(mesh.matrix);
    if (v.y > 0) {
      const r = Math.hypot(v.x, v.z);
      const t = THREE.MathUtils.smoothstep(v.y, 2, 9);   // 樹幹下段和樹根不動
      v.y += LIFT * r * t;
      v.x *= 1 + SPREAD * t;   // 向四面八方伸展
      v.z *= 1 + SPREAD * t;
    }
    pos.set([v.x, v.y, v.z], i * 3);
    top = Math.max(top, v.y);
  }
  const out = new THREE.BufferGeometry();
  out.setAttribute('position', new THREE.BufferAttribute(pos, 3));
  const uv = g.attributes.uv, uvs = new Float32Array(n * 2);
  for (let i = 0; i < n; i++) { uvs[i * 2] = uv.getX(i); uvs[i * 2 + 1] = uv.getY(i); }
  out.setAttribute('uv', new THREE.BufferAttribute(uvs, 2));
  out.setIndex(g.index);
  out.computeVertexNormals();
  g.dispose();
  mesh.geometry = out;
  mesh.position.set(0, 0, 0); mesh.quaternion.identity(); mesh.scale.set(1, 1, 1);
  return top;
}

const WIND = `
  uniform float uTime;
  uniform float uTreeH;
  vec3 windSway(vec3 p, float flutter) {
    float h = max(p.y, 0.) / uTreeH;
    float bend = h * h;
    // 整棵樹順風輕擺，再疊上樹梢各處不同步的小擺動
    float gust = sin(uTime * .9) * .6 + sin(uTime * 2.1 + 1.3) * .25 + .55;
    float local = sin(uTime * 1.7 + p.x * .35 + p.z * .27) * .35;
    vec3 d = vec3(1., 0., .35) * bend * (gust + local) * .55;
    // 葉子：每片各自細碎抖動
    float ph = dot(p, vec3(12.9898, 78.233, 37.719));
    d += flutter * vec3(sin(uTime * 7.3 + ph), sin(uTime * 9.1 + ph * 1.3) * .6, cos(uTime * 6.7 + ph)) * .09 * (.3 + h);
    return p + d;
  }`;

function addWind(material, uTime, flutter, treeH) {
  material.onBeforeCompile = shader => {
    shader.uniforms.uTime = uTime;
    shader.uniforms.uTreeH = treeH;
    shader.vertexShader = shader.vertexShader
      .replace('#include <common>', '#include <common>\n' + WIND)
      .replace('#include <begin_vertex>', `vec3 transformed = windSway(position, ${flutter.toFixed(1)});`);
  };
}

/**
 * 在 (x, z) 種一棵菩提樹，height 是地面以上的高度（場景單位）。
 * 回傳 { group, update(t, dt), dispose() }。
 */
export function plantBodhiTree(scene, { x, z, height, yaw = 0, uTime }) {
  const s = height / MODEL_HEIGHT;
  const group = new THREE.Group();
  group.position.set(x, 0, z);
  group.rotation.y = yaw;
  scene.add(group);

  const tl = new THREE.TextureLoader();
  const bark = tl.load(BASE + 'bark.jpg');
  bark.colorSpace = THREE.SRGBColorSpace;
  bark.wrapS = bark.wrapT = THREE.RepeatWrapping;
  const leaf = tl.load(BASE + 'leaf.png');
  leaf.colorSpace = THREE.SRGBColorSpace;

  const barkMat = new THREE.MeshStandardMaterial({ map: bark, color: '#cdbfb2', roughness: 0.95 });
  const leafMat = new THREE.MeshStandardMaterial({ map: leaf, alphaTest: 0.5, side: THREE.DoubleSide, roughness: 0.75, color: '#f4ffe0', emissive: '#2a4018' });   // 新綠
  const treeH = { value: MODEL_HEIGHT };   // 風擺的高度基準：載入後改成抬高後的樹高
  addWind(barkMat, uTime, 0, treeH);
  addWind(leafMat, uTime, 1, treeH);

  const tree = new THREE.Group();
  tree.scale.setScalar(s);
  group.add(tree);
  const loader = new GLTFLoader().setMeshoptDecoder(MeshoptDecoder);
  loader.load(BASE + 'bodhi-tree.glb', gltf => {
    let top = 0;
    gltf.scene.traverse(o => {
      if (!o.isMesh) return;
      top = Math.max(top, raiseBranches(o));
      o.material = o.name === 'leaves' ? leafMat : barkMat;
      o.frustumCulled = false;   // 隨風擺動後的範圍超出原本的包圍盒
    });
    // 葉子加密：同一層葉子再複製兩份，繞樹幹轉個角度、稍微縮放，樹冠更茂盛
    const leaves = gltf.scene.getObjectByName('leaves');
    if (leaves) {
      for (const [turn, k] of [[0.7, 0.93], [-1.3, 1.04]]) {
        const extra = new THREE.Mesh(leaves.geometry, leafMat);
        extra.rotation.y = turn;
        extra.scale.setScalar(k);
        extra.frustumCulled = false;
        leaves.parent.add(extra);
      }
    }
    // 樹枝抬高後整棵樹變高，縮回指定的高度
    treeH.value = top;
    tree.scale.setScalar(height / top);
    tree.add(gltf.scene);
  });

  // 飄落的菩提葉：在樹冠範圍內出生，順風打轉飄下，落地前縮小消失，再回到樹冠
  const leafSize = Math.max(0.5 * s * 1.6, height * 0.04);   // 比樹上的葉子稍大，遠遠也看得到
  const fallMat = new THREE.MeshStandardMaterial({ map: leaf, alphaTest: 0.5, side: THREE.DoubleSide, roughness: 0.8, color: '#e6d9a8' });
  const fall = new THREE.InstancedMesh(new THREE.PlaneGeometry(leafSize * 0.67, leafSize), fallMat, FALLING);
  fall.frustumCulled = false;
  group.add(fall);
  const crownR = 14 * s, crownLow = height * 0.45, crownHigh = height * 0.9;
  const leaves = Array.from({ length: FALLING }, () => spawn({}, true));
  function spawn(l, anywhere) {
    const a = Math.random() * Math.PI * 2, r = Math.sqrt(Math.random()) * crownR;
    l.x = Math.cos(a) * r;
    l.z = Math.sin(a) * r;
    l.y = anywhere ? Math.random() * crownHigh : crownLow + Math.random() * (crownHigh - crownLow);
    l.speed = (0.35 + Math.random() * 0.35) * height / 7;
    l.phase = Math.random() * 100;
    l.spin = new THREE.Vector3(Math.random(), Math.random(), Math.random()).normalize();
    l.spinRate = 1 + Math.random() * 2.5;
    return l;
  }
  const m = new THREE.Matrix4(), q = new THREE.Quaternion(), p = new THREE.Vector3(), sc = new THREE.Vector3();
  const tint = new THREE.Color();
  for (let i = 0; i < FALLING; i++) fall.setColorAt(i, tint.setHSL(0.15 - Math.random() * 0.08, 0.45, 0.55 + Math.random() * 0.2));

  return {
    group,
    update(t, dt) {
      const step = Math.min(dt, 0.1);
      for (let i = 0; i < FALLING; i++) {
        const l = leaves[i];
        // 下落、順風（+x）漂移、左右飄盪
        l.y -= l.speed * step;
        l.x += (0.25 + Math.sin(t * 0.9) * 0.15) * l.speed * step + Math.sin(t * 1.3 + l.phase) * 0.15 * s * step * 10;
        l.z += Math.cos(t * 1.1 + l.phase) * 0.1 * s * step * 10;
        if (l.y < 0) spawn(l, false);
        q.setFromAxisAngle(l.spin, t * l.spinRate + l.phase);
        const k = Math.min(1, l.y / (0.6 * s * 10));   // 接近地面時縮小，像是落進草叢
        p.set(l.x, l.y, l.z);
        sc.setScalar(Math.max(k, 0.001));
        fall.setMatrixAt(i, m.compose(p, q, sc));
      }
      fall.instanceMatrix.needsUpdate = true;
    },
    dispose() {
      bark.dispose(); leaf.dispose();
      barkMat.dispose(); leafMat.dispose(); fallMat.dispose();
      fall.geometry.dispose();
      tree.traverse(o => o.geometry?.dispose());
    },
  };
}
