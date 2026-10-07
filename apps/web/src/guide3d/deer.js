// 草原上的梅花鹿：Geodown 用 Meshy 生成的鹿（四足骨架、沒有附動畫），動作全部用程式擺骨頭：
// 四條腿依序交替的走路步態、低頭吃草、抬頭張望、甩尾巴。幾隻鹿各自在草地上隨機挑地方走過去，
// 走到了就停下來吃草或張望，再換下一個地方；路線避開河道、菩提樹和 Sunny。
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { MeshoptDecoder } from 'three/addons/libs/meshopt_decoder.module.js';
import { clone as cloneSkinned } from 'three/addons/utils/SkeletonUtils.js';

const MODEL = '/models/deer/deer.glb';   // 模型面向 +z，單位是公尺（含鹿角高 1.7）
const STRIDE = 1.25;      // 一個完整步態週期前進的距離（公尺，未縮放）
const SPEED = 0.75;       // 走路速度（公尺／秒，未縮放）
const TURN = 1.6;         // 轉身速度（弧度／秒）

// 骨頭名稱（Meshy SmartRig）：腿由上到下
const LEGS = [
  // [上段, 膝／跗關節, 蹄上方, 步態相位, 膝蓋彎的方向]：後腿的跗關節往前收、前腿的膝蓋往後收
  ['Bone_024', 'Bone_022', 'Bone_021', 0.0, -1],   // 右後
  ['Bone_045', 'Bone_043', 'Bone_042', 0.25, 1],   // 右前
  ['Bone_017', 'Bone_015', 'Bone_014', 0.5, -1],   // 左後
  ['Bone_037', 'Bone_035', 'Bone_034', 0.75, 1],   // 左前
];
const NECK = [['Bone_029', 0.55], ['Bone_028', 0.4], ['Bone_027', 0.3], ['Bone_026', 0.25], ['Bone_025', 0.25]];   // 低頭時各節彎多少
const CHEST = 'Bone_006', TAIL = ['Bone_005', 'Bone_004'], HIP = 'Bone_000';

const ease = (cur, target, rate, dt) => cur + (target - cur) * (1 - Math.exp(-rate * dt));
const rand = (a, b) => a + Math.random() * (b - a);

/**
 * 在場景裡放 count 隻鹿。
 * walkable(x, z)：這個地面位置能不能走（草地、不在河裡、不撞到東西）
 * pick()：隨機挑一個可以走去的地面位置 { x, z }
 * spawn(i)：第 i 隻鹿出現的位置（沒給就用 pick）；出現後很快就開始走
 * roam：每次走多遠（[最近, 最遠]，場景單位）
 * waterAt(x, z)：這裡有沒有水（0～1）；走進河裡時身體往下沉 wade，像涉水
 * 回傳 { update(t, dt), dispose() }
 */
export function addDeerHerd(scene, { count = 3, scale = 0.5, walkable, pick, spawn, roam = [1, 7], waterAt, wade = 0 }) {
  const herd = [];
  const shadowMat = new THREE.ShaderMaterial({
    transparent: true, depthWrite: false,
    vertexShader: 'varying vec2 vUv; void main(){ vUv = uv; gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.); }',
    fragmentShader: 'varying vec2 vUv; void main(){ float r = length(vUv - .5) * 2.; gl_FragColor = vec4(0.12, 0.1, 0.08, smoothstep(1., .2, r) * .32); }',
  });
  const shadowGeo = new THREE.PlaneGeometry(0.75, 1.7);
  let disposed = false, source = null, pose = null;

  new GLTFLoader().setMeshoptDecoder(MeshoptDecoder).load(MODEL, gltf => {
    if (disposed) return;
    source = gltf.scene;
    // 每根要動的骨頭：記下靜止姿勢，以及「鹿身體的橫軸、直軸」在牠父骨頭座標裡的方向
    source.updateMatrixWorld(true);
    const rig = {};
    const q = new THREE.Quaternion();
    source.traverse(o => {
      if (!o.isBone) return;
      o.parent.getWorldQuaternion(q).invert();
      rig[o.name] = {
        rest: o.quaternion.clone(),
        side: new THREE.Vector3(1, 0, 0).applyQuaternion(q),   // 繞橫軸：前後擺
        up: new THREE.Vector3(0, 1, 0).applyQuaternion(q),     // 繞直軸：左右轉
      };
    });
    source.traverse(o => {
      if (!o.isMesh) return;
      o.material.roughness = 0.9;
      o.frustumCulled = false;   // 骨頭動了以後超出原本的包圍盒
    });

    for (let i = 0; i < count; i++) {
      const s = scale * rand(0.88, 1.08);   // 大小略有不同
      const model = cloneSkinned(source);
      model.scale.setScalar(s);
      const group = new THREE.Group();
      group.add(model);
      const shadow = new THREE.Mesh(shadowGeo, shadowMat);
      shadow.rotation.x = -Math.PI / 2;
      shadow.position.y = 0.012;
      shadow.scale.setScalar(s);
      group.add(shadow);
      const start = spawn?.(i) ?? pick() ?? { x: 0, z: 0 };
      group.position.set(start.x, 0, start.z);
      const heading = rand(0, Math.PI * 2);
      group.rotation.y = heading;
      scene.add(group);
      const bones = {};
      model.traverse(o => { if (o.isBone) bones[o.name] = o; });
      const deer = {
        group, model, bones, s, heading,
        state: 'graze', timer: spawn ? rand(0.3, 2) : rand(1, 6),   // 一開始先低頭吃草，錯開出發時間
        wet: 0, shadow,
        target: null, phase: Math.random(), walk: 0, graze: 1, look: 0, lookYaw: 0, tail: 0, tailT: rand(2, 6),
        seed: Math.random() * 100,
      };
      herd.push(deer);
    }

    // 擺一根骨頭：靜止姿勢再繞身體的橫軸（或直軸）轉 angle
    const tmp = new THREE.Quaternion(), tmp2 = new THREE.Quaternion();
    pose = (deer, name, pitch, yaw = 0) => {
      const b = deer.bones[name], r = rig[name];
      if (!b || !r) return;
      tmp.setFromAxisAngle(r.side, pitch);
      if (yaw) tmp.premultiply(tmp2.setFromAxisAngle(r.up, yaw));
      b.quaternion.copy(tmp).multiply(r.rest);
    };
  });

  if (new URLSearchParams(location.search).has('debug')) window.__deer = herd;   // 除錯用

  /** 從 (x0, z0) 直走到 (x1, z1) 的路上都能走嗎 */
  function clearPath(x0, z0, x1, z1, self) {
    const n = Math.ceil(Math.hypot(x1 - x0, z1 - z0) / 0.3);
    for (let k = 1; k <= n; k++) {
      const x = x0 + (x1 - x0) * k / n, z = z0 + (z1 - z0) * k / n;
      if (!walkable(x, z)) return false;
    }
    // 不要和別的鹿挑到同一個地方
    return herd.every(o => o === self || !o.target || Math.hypot(o.target.x - x1, o.target.z - z1) > 1.2);
  }

  function chooseTarget(deer) {
    const { x, z } = deer.group.position;
    for (let k = 0; k < 20; k++) {
      const p = pick();
      if (!p) continue;
      const d = Math.hypot(p.x - x, p.z - z);
      if (d > roam[0] && d < roam[1] && clearPath(x, z, p.x, p.z, deer)) return p;
    }
    return null;
  }

  return {
    update(t, dt) {
      if (!pose) return;
      const step = Math.min(dt, 0.1);
      for (const deer of herd) {
        const g = deer.group;
        // 狀態：吃草／張望 → 時間到就挑下一個地方走過去 → 走到了再停下來
        deer.timer -= step;
        if (deer.state !== 'walk' && deer.timer <= 0) {
          deer.target = chooseTarget(deer);
          if (deer.target) deer.state = 'walk';
          else { deer.state = 'look'; deer.timer = rand(1.5, 3); }
        }
        let moving = 0;
        if (deer.state === 'walk') {
          const dx = deer.target.x - g.position.x, dz = deer.target.z - g.position.z;
          const dist = Math.hypot(dx, dz);
          if (dist < 0.08) {
            deer.target = null;
            deer.state = Math.random() < 0.65 ? 'graze' : 'look';
            deer.timer = deer.state === 'graze' ? rand(4, 10) : rand(2, 4);
          } else {
            // 一邊走一邊轉向目標，角度差大時放慢
            const want = Math.atan2(dx, dz);
            let diff = Math.atan2(Math.sin(want - deer.heading), Math.cos(want - deer.heading));
            deer.heading += THREE.MathUtils.clamp(diff, -TURN * step, TURN * step);
            moving = Math.max(0, Math.cos(diff)) * Math.min(1, dist / 0.4 + 0.3);
            const v = SPEED * deer.s * moving * step;
            const nx = g.position.x + Math.sin(deer.heading) * v, nz = g.position.z + Math.cos(deer.heading) * v;
            if (walkable(nx, nz)) { g.position.x = nx; g.position.z = nz; }
            else { deer.target = null; deer.state = 'look'; deer.timer = rand(1, 2); }
          }
        }
        g.rotation.y = deer.heading;
        // 涉水：走進河裡時慢慢沉下去，水面蓋住小腿；水裡不畫地上的影子
        if (waterAt) {
          deer.wet = ease(deer.wet, waterAt(g.position.x, g.position.z) > 0.3 ? 1 : 0, 3, step);
          g.position.y = -wade * deer.wet;
          deer.shadow.visible = deer.wet < 0.5;
        }

        // 動作的權重慢慢過渡，切換狀態時不會突然跳動
        deer.walk = ease(deer.walk, deer.state === 'walk' ? Math.max(moving, 0.35) : 0, 6, step);
        deer.graze = ease(deer.graze, deer.state === 'graze' ? 1 : 0, 2.2, step);
        deer.look = ease(deer.look, deer.state === 'look' ? 1 : 0, 3, step);
        deer.phase += SPEED * moving * step / STRIDE;

        const w = deer.walk, P = deer.phase * Math.PI * 2;
        // 腿：上段前後擺，往前跨（擺動期）時膝蓋收起、蹄抬離地面
        for (const [upper, knee, pastern, off, dir] of LEGS) {
          const a = P + off * Math.PI * 2;
          const swing = Math.max(0, -Math.cos(a));   // 往前跨的那半個週期
          pose(deer, upper, Math.sin(a) * 0.32 * w);
          pose(deer, knee, dir * swing * 0.75 * w);
          pose(deer, pastern, dir * swing * 0.35 * w);
        }
        // 身體隨步伐輕輕上下起伏
        deer.model.position.y = -Math.abs(Math.sin(P * 2)) * 0.015 * w;
        pose(deer, HIP, Math.sin(P * 2) * 0.02 * w);
        // 頭：吃草時整條脖子往下彎；張望時抬頭左右轉；走路時隨步伐點頭
        const lookYaw = Math.sin(t * 0.5 + deer.seed) * 0.6 * deer.look;
        const chew = Math.sin(t * 5 + deer.seed) * 0.04 * deer.graze;
        pose(deer, CHEST, 0.18 * deer.graze);
        for (const [name, k] of NECK) {
          const nod = Math.sin(P * 2 + 0.6) * 0.04 * w;
          pose(deer, name, k * (1.25 * deer.graze - 0.15 * deer.look) + nod + chew, name === 'Bone_027' || name === 'Bone_026' ? lookYaw * 0.5 : 0);
        }
        // 尾巴：偶爾甩一甩
        deer.tailT -= step;
        if (deer.tailT <= 0) { deer.tail = 1; deer.tailT = rand(2, 7); }
        deer.tail = Math.max(0, deer.tail - step * 1.5);
        const flick = Math.sin(t * 18) * 0.5 * deer.tail;
        pose(deer, TAIL[0], -0.2 * deer.tail, flick);
        pose(deer, TAIL[1], 0, flick);
      }
    },
    dispose() {
      disposed = true;
      for (const d of herd) scene.remove(d.group);
      source?.traverse(o => { o.geometry?.dispose(); o.material?.map?.dispose(); o.material?.dispose(); });
      shadowGeo.dispose(); shadowMat.dispose();
    },
  };
}
