// 松林雪山：薄霧松林間的一小片空地，遠方雪峰，Sunny 站在雪岸旁那條河的中間。
// 陰天的冷色散射光，沒有菩提樹。共用的全景圖做法見 scene-panorama.js。
import * as THREE from 'three';
import { buildPanoramaScene } from './scene-panorama.js';

export function buildPinesScene(scene, camera, renderer) {
  return buildPanoramaScene(scene, camera, renderer, {
    panorama: '/scenes/pines-snow.jpg',
    panorama8k: '/scenes/pines-snow-8k.jpg',
    riverMask: '/scenes/pines-snow-river.png',
    eyeHeight: 1.2,
    // 遠景構圖：雪峰和松林在上、河在下；左右避開近處的兩棵大樹幹
    viewU: 0.29,
    viewPitch: 6,        // 稍微抬頭，讓雪峰露出來（頁面又把畫面往上挪了一點）
    spot: [0.31, 0.548], // 河道裡，雪岸前方
    wideHFov: 66,
    background: '#c9d1d6',   // 全景圖載入前的陰天灰藍
    glint: [0.95, 0.98, 1],
    sunDir: new THREE.Vector3(0.2, 0.8, 0.6),    // 雲層後的天光，從上方偏前方灑下
    sun: ['#eef2f6', 1.2],
    hemi: ['#dfe7ee', '#6b5a48', 1.4],
    fill: 0.7,
    ripple: [0.96, 0.98, 1],
  });
}
