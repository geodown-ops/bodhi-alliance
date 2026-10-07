// 晨霧草原：Sunny 站在河流分岔處偏左的水裡，前方草地上一棵菩提樹隨風擺動、落葉紛飛。
// 其他物件（湖、蘆葦、魚、鳥）都不放。共用的全景圖做法見 scene-panorama.js。
import * as THREE from 'three';
import { buildPanoramaScene } from './scene-panorama.js';
import { plantBodhiTree } from './bodhi-tree.js';
import { addDeerHerd } from './deer.js';
import { addLotus } from './lotus.js';

// 菩提樹：種在畫面中間、比 Sunny 更靠近鏡頭的草地上（不遮到她）（全景圖座標 u, v）與地面以上的高度
const TREE_SPOT = [0.1, 0.565];
const TREE_HEIGHT = 4.5;
// 梅花鹿：從畫面兩邊的河裡出發，涉水走向 Sunny 身後的草地，之後在草地上大範圍散步
const DEER_COUNT = 3;
const DEER_SCALE = 0.5 * 2 / 3;   // 比原本小三分之一
const DEER_BEHIND = 0.8, DEER_FAR = 15, DEER_ANGLE = 0.85;   // 比 Sunny 遠 0.8 以上；左右到畫面邊緣
// 蓮花：只放三朵散在河面上（種類與數量）
const LOTUS = ['flower', 'serenity', 'pod'];
const LOTUS_AWAY = 0.5;   // 和 Sunny 在畫面上至少錯開這麼多（鏡頭看出去的角度，弧度）

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
    extras({ scene, uTime, at, riverAt, riverReady, standAt, camZ, q }) {
      const treeAt = at(Number(q.get('treeU')) || TREE_SPOT[0], Number(q.get('treeV')) || TREE_SPOT[1]);
      const tree = plantBodhiTree(scene, { ...treeAt, height: Number(q.get('treeH')) || TREE_HEIGHT, yaw: 0.6, uTime });
      // 鹿能走的地方：比 Sunny 離鏡頭更遠（在她身後）、左右不超出畫面，不撞到菩提樹和 Sunny
      const sunnyDist = Math.hypot(standAt.x, camZ - standAt.z), near = sunnyDist + DEER_BEHIND;
      const walkable = (x, z) => {
        const dz = camZ - z, d = Math.hypot(x, dz);
        if (d < near || d > DEER_FAR || Math.abs(Math.atan2(x, dz)) > DEER_ANGLE) return false;
        return Math.hypot(x - treeAt.x, z - treeAt.z) > 1 && Math.hypot(x - standAt.x, z - standAt.z) > 1.2;
      };
      const sunnyAngle = Math.atan2(standAt.x, camZ - standAt.z);
      const polar = (a, d) => ({ x: Math.sin(a) * d, z: camZ - Math.cos(a) * d });
      // 目的地只挑草地（路上可以涉水過河）
      const pick = () => {
        const p = polar((Math.random() * 2 - 1) * DEER_ANGLE, near + Math.random() * (DEER_FAR - near));
        return walkable(p.x, p.z) && riverAt(p.x, p.z) < 0.3 ? p : null;
      };
      // 出發點：輪流在畫面最左、最右兩側
      const spawn = i => {
        const side = i % 2 ? 1 : -1;
        for (let k = 0; k < 30; k++) {
          const p = polar(side * (DEER_ANGLE - 0.02 - Math.random() * 0.1), near + 0.5 + (i >> 1) * 3 + Math.random() * 2.5);   // 同一側的鹿前後錯開
          if (walkable(p.x, p.z)) return p;
        }
        return null;
      };
      const herd = addDeerHerd(scene, {
        count: Number(q.get('deer') ?? DEER_COUNT), scale: DEER_SCALE, walkable, pick, spawn,
        roam: [3, 14], waterAt: riverAt, wade: 0.09,
      });
      // 蓮花：在河道裡挑位置（固定的亂數，每次打開都一樣），遠離 Sunny，彼此也留點距離
      let lotus = null, disposed = false;
      riverReady.then(() => {
        if (disposed) return;
        let seed = 7;
        const rnd = () => (seed = (seed * 16807) % 2147483647) / 2147483647;
        const spots = [];
        for (let k = 0; k < 2000 && spots.length < LOTUS.length; k++) {
          const a = -0.8 + rnd() * 1.25, p = polar(a, 4 + rnd() * 5);   // 4 以內會被輸入框擋住；右邊太遠會被標題擋住
          const inWater = [[0, 0], [0.3, 0], [-0.3, 0], [0, 0.3], [0, -0.3]].every(([ox, oz]) => riverAt(p.x + ox, p.z + oz) > 0.6);
          if (!inWater || Math.abs(a - sunnyAngle) < LOTUS_AWAY) continue;
          if (spots.some(o => Math.hypot(o.x - p.x, o.z - p.z) < 1)) continue;
          spots.push({ ...p, kind: LOTUS[spots.length] });
        }
        lotus = addLotus(scene, spots);
      });
      return {
        update(t, dt) { tree.update(t, dt); herd.update(t, dt); lotus?.update(t); },
        dispose() { disposed = true; tree.dispose(); herd.dispose(); lotus?.dispose(); },
      };
    },
  });
}
