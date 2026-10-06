// 晨霧草原：Sunny 站在河流分岔處偏左的水裡，前方草地上一棵菩提樹隨風擺動、落葉紛飛。
// 其他物件（湖、蘆葦、魚、鳥）都不放。共用的全景圖做法見 scene-panorama.js。
import * as THREE from 'three';
import { buildPanoramaScene } from './scene-panorama.js';
import { plantBodhiTree } from './bodhi-tree.js';

// 菩提樹：種在畫面中間、比 Sunny 更靠近鏡頭的草地上（不遮到她）（全景圖座標 u, v）與地面以上的高度
const TREE_SPOT = [0.1, 0.565];
const TREE_HEIGHT = 4.5;

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
    extras({ scene, uTime, at, q }) {
      const treeAt = at(Number(q.get('treeU')) || TREE_SPOT[0], Number(q.get('treeV')) || TREE_SPOT[1]);
      return plantBodhiTree(scene, { ...treeAt, height: Number(q.get('treeH')) || TREE_HEIGHT, yaw: 0.6, uTime });
    },
  });
}
