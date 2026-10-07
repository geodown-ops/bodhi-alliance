// 河面上的蓮花：Geodown 用 Meshy 生成的蓮花、蓮蓬、睡蓮（Lotus Serenity）三種模型，少量散在河裡，
// 隨水波輕輕上下浮動、慢慢打轉。位置由場景給（只挑河道裡的地方）。
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { MeshoptDecoder } from 'three/addons/libs/meshopt_decoder.module.js';

const BASE = '/models/lotus/';
// 每種模型：檔名、在場景裡的寬度（場景單位；Sunny 約 0.8 高）、沉進水面的比例
const KINDS = {
  flower: { file: 'flower.glb', size: 0.32, sink: 0.25 },
  pod: { file: 'pod.glb', size: 0.12, sink: 0.35 },          // 蓮蓬：size 指寬度，連莖約高 0.25
  serenity: { file: 'serenity.glb', size: 0.5, sink: 0.3 },
};

/**
 * spots：[{ kind: 'flower' | 'pod' | 'serenity', x, z, scale? }]（scale：額外放大倍數）
 * 回傳 { update(t), dispose() }
 */
export function addLotus(scene, spots) {
  const group = new THREE.Group();
  scene.add(group);
  const items = [];
  const loaded = [];
  const loader = new GLTFLoader().setMeshoptDecoder(MeshoptDecoder);
  const box = new THREE.Box3(), size = new THREE.Vector3();
  let disposed = false;
  for (const [kind, k] of Object.entries(KINDS)) {
    const mine = spots.filter(s => s.kind === kind);
    if (!mine.length) continue;
    loader.load(BASE + k.file, gltf => {
      if (disposed) return;
      loaded.push(gltf.scene);
      box.setFromObject(gltf.scene).getSize(size);
      const scale = k.size / Math.max(size.x, size.z);
      for (const s of mine) {
        const o = gltf.scene.clone();
        o.scale.setScalar(scale * (s.scale ?? 1) * (0.85 + Math.random() * 0.3));
        const baseY = -box.min.y * o.scale.y - size.y * o.scale.y * k.sink;   // 底部沉進水裡一點
        o.position.set(s.x, baseY, s.z);
        o.rotation.y = Math.random() * Math.PI * 2;
        group.add(o);
        items.push({ o, baseY, yaw: o.rotation.y, seed: Math.random() * 100, spin: (Math.random() - 0.5) * 0.06 });
      }
    });
  }
  return {
    update(t) {
      for (const it of items) {
        it.o.position.y = it.baseY + Math.sin(t * 0.9 + it.seed) * 0.006;   // 隨水波上下浮動
        it.o.rotation.y = it.yaw + t * it.spin;
        it.o.rotation.z = Math.sin(t * 0.7 + it.seed) * 0.03;
      }
    },
    dispose() {
      disposed = true;
      scene.remove(group);
      for (const s of loaded) s.traverse(o => {
        o.geometry?.dispose();
        for (const m of [o.material].flat()) if (m) { for (const v of Object.values(m)) if (v?.isTexture) v.dispose(); m.dispose(); }
      });
    },
  };
}
