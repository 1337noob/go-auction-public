import http from 'k6/http';
import { check, sleep } from 'k6';

/**
 * k6 script to place bids against pre-existing auctions.
 * Pass auction IDs via env AUCTION_IDS (comma-separated).
 * Example:
 *   AUCTION_IDS=id1,id2 BID_BASE=1000 BID_STEP=10 k6 run auction-load.js
 */

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081';
const AUCTION_IDS = (__ENV.AUCTION_IDS || '')
  .split(',')
  .map((s) => s.trim())
  .filter((s) => s.length > 0);
const BID_BASE = Number(__ENV.BID_BASE || 1000);
const BID_STEP = Number(__ENV.BID_STEP || 10);

export const options = {
  vus: Number(__ENV.VUS || 20),
  duration: __ENV.DURATION || '1m',
  thresholds: {
    http_req_failed: ['rate<0.01'], // keep failures under 1%
    http_req_duration: ['p(95)<500'], // 95% under 500ms
  },
};

function postJson(path, payload) {
  return http.post(`${BASE_URL}${path}`, JSON.stringify(payload), {
    headers: { 'Content-Type': 'application/json' },
    tags: { path },
  });
}

export function setup() {
  if (AUCTION_IDS.length === 0) {
    throw new Error('AUCTION_IDS env is required (comma-separated list)');
  }
  return { auctions: AUCTION_IDS };
}

export default function (data) {
  const auctions = data.auctions;
  if (!auctions || auctions.length === 0) {
    return;
  }

  const target = auctions[Math.floor(Math.random() * auctions.length)];
  const bump = Math.floor(Math.random() * 10) * BID_STEP;
  const body = {
    user_id: `user-${__VU}-${Math.floor(Math.random() * 1000)}`,
    amount: BID_BASE + BID_STEP + bump,
  };

  const res = postJson(`/auctions/${target}/bids`, body);
  check(res, {
    'bid accepted (202/200)': (r) => r.status === 202 || r.status === 200,
    'bid not client error': (r) => r.status < 400,
  });

  sleep(1);
}

