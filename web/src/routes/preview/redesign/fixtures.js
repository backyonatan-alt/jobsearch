// Fixture dataset from the Claude Design handoff — identical across all
// redesign preview pages. 13 active / 3 closed.
export const logo = (d) => `https://www.google.com/s2/favicons?domain=${d}&sz=64`;
export const logoLg = (d) => `https://www.google.com/s2/favicons?domain=${d}&sz=128`;

export const CATS = [
  { count: 1, label: 'to prep',   color: '#e0641f', labelColor: '#16181c', note: 'panel tomorrow',    domains: ['kayma.com'] },
  { count: 1, label: 'to decide', color: '#16a34a', labelColor: '#16181c', note: 'offer · $210k',     domains: ['linear.app'] },
  { count: 2, label: 'to nudge',  color: '#b3372a', labelColor: '#16181c', note: 'quiet 3+ weeks',    domains: ['hibob.com', 'rylo.com'] },
  { count: 9, label: 'waiting',   color: '#16181c', labelColor: '#6f7680', note: 'nothing to do yet', domains: ['wiz.io', 'fiverr.com', 'gong.io'], extra: 6 },
  { count: 3, label: 'closed',    color: '#9aa1ab', labelColor: '#6f7680', note: 'kept for the record', domains: ['wix.com', 'deel.com', 'papayaglobal.com'], gray: true }
];

const act = (company, d, meta, hot, hotColor, border, btn, btnBg, btnColor, btnBorder) => ({
  kind: 'act', company, meta, hot, hotColor, domain: d, border, btn, btnBg, btnColor, btnBorder
});
const quiet = (company, d, meta, quietLabel) => ({ kind: 'quiet', company, meta, domain: d, quiet: quietLabel });
const exit = (company, d, meta) => ({ kind: 'exit', company, meta, domain: d, btn: 'Reopen' });

export const HOME_ROWS = [
  act('Kayma', 'kayma.com', 'VP Product · panel', 'tomorrow 10:00', '#c05310', '#cdddfb', 'Prep now →', '#2463eb', '#fff', '#2463eb'),
  act('Linear', 'linear.app', 'Senior PM · offer', '$210k · answer this week', '#1d7a4f', '#cfe5d2', 'Decide →', '#16a34a', '#fff', '#16a34a'),
  act('HiBob', 'hibob.com', 'Senior PM ·', 'quiet 21 days', '#b3372a', '#f2d4cf', 'Follow up', '#fff', '#2463eb', '#cdddfb'),
  act('Rylo', 'rylo.com', 'Product Manager ·', 'quiet 27 days', '#b3372a', '#f2d4cf', 'Follow up', '#fff', '#2463eb', '#cdddfb'),
  quiet('Wiz', 'wiz.io', 'Product Lead · screen booked Thu 15:00', 'on track'),
  quiet('Melio', 'meliopayments.com', 'Senior PM · recruiter is scheduling', 'waiting'),
  quiet('Fiverr', 'fiverr.com', 'Senior PM · applied today via LinkedIn', 'too early'),
  quiet('Gong', 'gong.io', 'Group PM · applied 5d ago via referral', 'waiting'),
  quiet('Lemonade', 'lemonade.com', 'Senior PM · followed up once · 40d', 'waiting'),
  quiet('Astelia', 'astelia.com', 'Product Manager · applied 40d ago', 'waiting'),
  quiet('JFrog', 'jfrog.com', 'Senior PM · followed up once · 47d', 'waiting'),
  quiet('Eleos Health', 'eleos.health', 'Senior PM · followed up once · 48d', 'waiting'),
  quiet('Riverside', 'riverside.fm', 'Senior PM · saved to wishlist 2d ago', 'apply when ready'),
  exit('Wix', 'wix.com', 'Senior PM · rejected after screen · 12d'),
  exit('Deel', 'deel.com', 'Senior PM · position closed · 8d'),
  exit('Papaya Global', 'papayaglobal.com', 'Product Manager · withdrawn · 30d')
];

const row = (company, d, meta, hot, hotColor, action) => ({
  kind: 'row', company, meta, hot: hot || '', hotColor: hotColor || '#8a9099', action: action || null, domain: d
});
const exitRow = (company, d, meta) => ({ kind: 'exit', company, meta, hot: '', action: 'Reopen', domain: d });

export const APP_GROUPS = [
  { name: 'Interview', count: 1, labelColor: '#e0641f', note: '', noteColor: '#8a9099', bg: '#fff', rows: [
    row('Kayma', 'kayma.com', 'VP Product ·', 'panel tomorrow 10:00', '#c05310', 'Prep →')
  ]},
  { name: 'Offer', count: 1, labelColor: '#1d7a4f', note: '', noteColor: '#8a9099', bg: '#fff', rows: [
    row('Linear', 'linear.app', 'Senior PM ·', '$210k · answer this week', '#1d7a4f', 'Decide →')
  ]},
  { name: 'Screen', count: 2, labelColor: '#0e9f6e', note: '', noteColor: '#8a9099', bg: '#fff', rows: [
    row('Wiz', 'wiz.io', 'Product Lead · booked Thu 15:00'),
    row('Melio', 'meliopayments.com', 'Senior PM · recruiter is scheduling')
  ]},
  { name: 'Applied', count: 8, labelColor: '#8a9099', note: '— 2 quiet too long', noteColor: '#b3372a', bg: '#fff', rows: [
    row('HiBob', 'hibob.com', 'Senior PM ·', 'quiet 21d', '#b3372a', 'Follow up'),
    row('Rylo', 'rylo.com', 'Product Manager ·', 'quiet 27d', '#b3372a', 'Follow up'),
    row('Fiverr', 'fiverr.com', 'Senior PM · applied today via LinkedIn'),
    row('Gong', 'gong.io', 'Group PM · applied 5d ago via referral'),
    row('Lemonade', 'lemonade.com', 'Senior PM · followed up once · 40d'),
    row('Astelia', 'astelia.com', 'Product Manager · applied 40d ago'),
    row('JFrog', 'jfrog.com', 'Senior PM · followed up once · 47d'),
    row('Eleos Health', 'eleos.health', 'Senior PM · followed up once · 48d')
  ]},
  { name: 'Wishlist', count: 1, labelColor: '#8a9099', note: '', noteColor: '#8a9099', bg: '#fff', rows: [
    row('Riverside', 'riverside.fm', 'Senior PM · saved 2d ago ·', '', '#8a9099', 'Apply →')
  ]},
  { name: 'No longer in play', count: 3, labelColor: '#b8bdc4', note: '— kept forever, reopen if the req comes back', noteColor: '#b8bdc4', bg: '#fbfbf9', rows: [
    exitRow('Wix', 'wix.com', 'Senior PM · rejected after screen · 12d'),
    exitRow('Deel', 'deel.com', 'Senior PM · position closed · 8d'),
    exitRow('Papaya Global', 'papayaglobal.com', 'Product Manager · withdrawn · 30d')
  ]}
];
