// 用 solc-js 編譯 BodhiCoin.sol，把 ABI 與部署用 bytecode 寫到 server/internal/chain/，
// api 服務會內嵌它們，第一次啟用鏈上功能時自己部署合約。
import { readFileSync, writeFileSync } from 'node:fs';
import { createRequire } from 'node:module';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import solc from 'solc';

const here = dirname(fileURLToPath(import.meta.url));
const require = createRequire(import.meta.url);
const out = join(here, '../server/internal/chain');

const input = {
  language: 'Solidity',
  sources: { 'BodhiCoin.sol': { content: readFileSync(join(here, 'BodhiCoin.sol'), 'utf8') } },
  settings: {
    optimizer: { enabled: true, runs: 200 },
    evmVersion: 'paris',
    outputSelection: { '*': { '*': ['abi', 'evm.bytecode.object'] } },
  },
};

function findImports(path) {
  try {
    return { contents: readFileSync(require.resolve(path), 'utf8') };
  } catch {
    return { error: `not found: ${path}` };
  }
}

const result = JSON.parse(solc.compile(JSON.stringify(input), { import: findImports }));
const errors = (result.errors ?? []).filter((e) => e.severity === 'error');
if (errors.length) {
  console.error(errors.map((e) => e.formattedMessage).join('\n'));
  process.exit(1);
}
const c = result.contracts['BodhiCoin.sol'].BodhiCoin;
writeFileSync(join(out, 'BodhiCoin.abi.json'), JSON.stringify(c.abi, null, 2) + '\n');
writeFileSync(join(out, 'BodhiCoin.bin'), c.evm.bytecode.object + '\n');
console.log(`solc ${solc.version()}: wrote ${out}/BodhiCoin.{abi.json,bin}`);
