// 晨霧草原：Sunny 站在河流分岔處偏左的水裡，前方草地上一棵菩提樹隨風擺動、落葉紛飛。
// 其他物件（湖、蘆葦、魚、鳥）都不放。共用的全景圖做法見 scene-panorama.js。
import * as THREE from 'three';
import { buildPanoramaScene } from './scene-panorama.js';
import { plantBodhiTree } from './bodhi-tree.js';
import { addDeerHerd } from './deer.js';

// 菩提樹：種在畫面中間、比 Sunny 更靠近鏡頭的草地上（不遮到她）（全景圖座標 u, v）與地面以上的高度
const TREE_SPOT = [0.1, 0.565];
const TREE_HEIGHT = 4.5;
// 梅花鹿：在鏡頭前方這片草地上隨機走動（離鏡頭的距離範圍、左右可走的角度）
const DEER_COUNT = 3;
const DEER_NEAR = 4.8, DEER_FAR = 12, DEER_ANGLE = 0.7;   // 不走到鏡頭前太近，免得擋住標題

export function buildMeadowScene(scene, camera, renderer) {
  return buildPanoramaScene(scene, camera, renderer, {
    panorama: '/scenes/meadow-dawn.jpg',
    panorama8k: '/scenes/meadow-dawn-8k.jpg',
    riverMask: '/scenes/meadow-dawn-river.png',
    eyeHeight: 1.2,
    // 遠景構圖：比照 Geodown 在 Skybox 網站上截的角度（山坡上的樹、河流在草地前分岔）
    viewU: 0.078,
    viewPitch: 2,        // 頁面把畫面往上挪了一些（setViewOffset），這裡稍微抬頭補回來
    spot: [0.043, 0.549],
    wideHFov: 66,        // 和 Skybox 網站上的畫面一樣寬
    background: '#d9a6b8',   // 全景圖載入前的粉色晨空
    glint: [1, 0.96, 0.98],
    sunDir: new THREE.Vector3(-0.3, 0.5, 0.8),   // 晨光從觀眾這一側照到臉上
    sun: ['#ffe2d0', 1.9],
    hemi: ['#f3d6e4', '#7d9458', 1.2],
    fill: 0.6,
    ripple: [1, 0.96, 0.97],
    extras({ scene, uTime, at, riverAt, standAt, camZ, q }) {
      const treeAt = at(Number(q.get('treeU')) || TREE_SPOT[0], Number(q.get('treeV')) || TREE_SPOT[1]);
      const tree = plantBodhiTree(scene, { ...treeAt, height: Number(q.get('treeH')) || TREE_HEIGHT, yaw: 0.6, uTime });
      // 草地上能走的地方：在鏡頭前方的範圍內、離河岸留一點距離、不撞到菩提樹和 Sunny
      const walkable = (x, z) => {
        const dz = camZ - z, d = Math.hypot(x, dz);
        if (d < DEER_NEAR || d > DEER_FAR || Math.abs(Math.atan2(x, dz)) > DEER_ANGLE) return false;
        if (Math.hypot(x - treeAt.x, z - treeAt.z) < 1 || Math.hypot(x - standAt.x, z - standAt.z) < 1.8) return false;
        for (const [ox, oz] of [[0, 0], [0.3, 0], [-0.3, 0], [0, 0.3], [0, -0.3]]) if (riverAt(x + ox, z + oz) > 0.1) return false;
        return true;
      };
      const pick = () => {
        const a = (Math.random() * 2 - 1) * DEER_ANGLE, d = DEER_NEAR + Math.random() * (DEER_FAR - DEER_NEAR);
        const p = { x: Math.sin(a) * d, z: camZ - Math.cos(a) * d };
        return walkable(p.x, p.z) ? p : null;
      };
      const herd = addDeerHerd(scene, { count: Number(q.get('deer') ?? DEER_COUNT), scale: 0.5, walkable, pick });
      return {
        update(t, dt) { tree.update(t, dt); herd.update(t, dt); },
        dispose() { tree.dispose(); herd.dispose(); },
      };
    },
  });
}
