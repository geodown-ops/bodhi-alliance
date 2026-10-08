// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/// 菩提幣（BODHI）。照 SPEC：小數 2 位、總量固定 5 億，部署時全部鑄給金庫。
/// 幣只能在「機構地址」（金庫、中心、共好企業、平台營運地址）與其他地址之間流動，
/// 兩邊都不是機構地址的轉帳（會員對會員）一律拒絕，所以會員之間不能互轉。
contract BodhiCoin is ERC20, Ownable {
    uint256 public constant TOTAL_SUPPLY = 500_000_000 * 100;

    mapping(address => bool) public institutional;

    event InstitutionalSet(address indexed account, bool allowed);

    error PersonalTransferNotAllowed(address from, address to);

    constructor(address treasury) ERC20(unicode"菩提幣", "BODHI") Ownable(msg.sender) {
        institutional[treasury] = true;
        emit InstitutionalSet(treasury, true);
        if (msg.sender != treasury) {
            institutional[msg.sender] = true;
            emit InstitutionalSet(msg.sender, true);
        }
        _mint(treasury, TOTAL_SUPPLY);
    }

    function decimals() public pure override returns (uint8) {
        return 2;
    }

    /// 決策小組通過後，由營運地址把中心、共好企業等加入或移出機構名單。
    function setInstitutional(address account, bool allowed) external onlyOwner {
        institutional[account] = allowed;
        emit InstitutionalSet(account, allowed);
    }

    function _update(address from, address to, uint256 value) internal override {
        if (from != address(0) && to != address(0) && !institutional[from] && !institutional[to]) {
            revert PersonalTransferNotAllowed(from, to);
        }
        super._update(from, to, value);
    }
}
