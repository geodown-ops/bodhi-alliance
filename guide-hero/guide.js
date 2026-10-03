// 首頁覺行小組線上組長：載入 3D 角色，處理單一輸入框的提問，並把回答以字幕顯示。
import { createAvatar } from './avatar.js';

const API = (window.BODHI_GUIDE_API || '') + '/api/chat';
const MODEL_URL = window.BODHI_GUIDE_MODEL || 'models/bodhi.vrm';
const $ = id => document.getElementById(id);
const form = $('ask'), input = $('askInput'), btn = $('askBtn');
const caption = $('caption'), capQ = $('capQ'), capA = $('capA');

// 窄螢幕放不下完整提示，改用短版（仍保留「請勿輸入個資」）
const narrow = matchMedia('(max-width: 480px)');
const setPlaceholder = () => { input.placeholder = narrow.matches ? '向覺行小組線上組長提問（請勿輸入個資）' : '向覺行小組線上組長提問（AI 回答・請勿輸入個人資料）'; };
narrow.addEventListener('change', setPlaceholder);
setPlaceholder();

// 角色載入失敗（不支援 WebGL、網路中斷）時，問答照常運作，只是沒有 3D 角色
let avatar = { setState() {}, speak() {}, finish() {} };
createAvatar($('avatar'), MODEL_URL, { onProgress: p => { $('bar').style.width = (p * 100).toFixed(0) + '%'; } })
  .then(a => { avatar = a; $('loading').classList.add('done'); })
  .catch(err => { console.error('avatar', err); $('loadText').textContent = ''; $('barWrap').remove(); });

const history = [];   // 對話只存在這個分頁，重新整理就清空
const MAX_TURNS = 20;

function show(q, a, cls = '') {
  caption.hidden = false;
  capQ.textContent = q;
  capA.className = 'g-a ' + cls;
  capA.textContent = a;
  caption.scrollTop = caption.scrollHeight;
}

async function ask(q) {
  q = q.trim(); if (!q) return;
  history.push({ role: 'user', content: q });
  while (history.length > MAX_TURNS) history.splice(0, 2);
  input.value = ''; btn.disabled = input.disabled = true;
  show(q, '');
  capA.innerHTML = '<span class="dots"></span>';
  avatar.setState('thinking');

  let answer = '';
  try {
    const res = await fetch(API, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ messages: history }) });
    if (!res.ok) throw new Error((await res.text()) || res.statusText);
    const reader = res.body.getReader(), dec = new TextDecoder();
    let buf = '';
    for (;;) {
      const { value, done } = await reader.read(); if (done) break;
      buf += dec.decode(value, { stream: true });
      let i;
      while ((i = buf.indexOf('\n\n')) >= 0) {
        const chunk = buf.slice(0, i); buf = buf.slice(i + 2);
        const ev = /^event: (.*)$/m.exec(chunk)?.[1], data = JSON.parse(/^data: (.*)$/m.exec(chunk)?.[1] || '{}');
        if (ev === 'delta') {
          answer += data.text;
          show(q, answer);
          avatar.speak(data.text);
        } else if (ev === 'error') throw new Error(data.message);
      }
    }
    if (answer) history.push({ role: 'assistant', content: answer.slice(0, 2000) });
    else history.pop();
    avatar.finish();
  } catch (e) {
    history.pop();
    show(q, e.message || '連線失敗，請稍後再試。', 'err');
    avatar.setState('idle');
  } finally {
    btn.disabled = input.disabled = false;
    input.focus({ preventScroll: true });
  }
}

form.addEventListener('submit', e => { e.preventDefault(); ask(input.value); });
