export type DocBlock =
  | { type: "p"; text: string }
  | { type: "ul"; items: string[] }
  | { type: "code"; text: string };

export type DocSection = {
  id: string;
  title: string;
  blocks: DocBlock[];
};

const zh: DocSection[] = [
  {
    id: "version",
    title: "版本号",
    blocks: [
      { type: "p", text: "侧边栏版本来自编译时写入的版本。使用 docker compose up -d --build 且没有设置 RESIN_VERSION 时，读取仓库根目录的 VERSION 文件。当前是 v1.3.0。" },
      { type: "p", text: "需要临时覆盖时，在构建环境设置 RESIN_VERSION。发布流水线仍使用 Git 标签作为版本。" },
    ],
  },
  {
    id: "subscriptions",
    title: "订阅",
    blocks: [
      { type: "p", text: "远程订阅支持单独的 User-Agent。留空时使用 clash.meta。本地订阅不能设置 User-Agent。" },
      { type: "p", text: "探测间隔留空时使用全局设置。填写后只缩短或指定该订阅节点的探测周期，最短 10 秒。一个节点属于多个订阅时，采用其中最短的正间隔。" },
      { type: "p", text: "批量导入按行读取订阅链接。空行和 # 开头的行会跳过。名称默认取域名中间段，也可以填写正则；有捕获组时用第一组。" },
      { type: "ul", items: ["纯文本代理行支持 IP:PORT 和 IP:PORT:USER:PASS。", "也支持 vmess、vless、trojan、ss、http、socks 以及 sing-box、Clash 订阅。"] },
    ],
  },
  {
    id: "healthy",
    title: "健康节点订阅",
    blocks: [
      { type: "p", text: "把当前健康且已有出口 IP 的节点导出成客户端订阅。管理接口使用管理员令牌，手机导入使用代理令牌。" },
      { type: "code", text: "GET /<代理令牌>/api/v1/healthy-subscription\nGET /<代理令牌>/api/v1/healthy-subscription?format=sing-box" },
      { type: "p", text: "默认返回 base64 编码的 URI 列表。format=sing-box 返回 sing-box JSON。" },
    ],
  },
  {
    id: "platforms",
    title: "平台与路由",
    blocks: [
      { type: "p", text: "打开「仅 IPv4」后，IPv6 出口、IPv6 服务器地址，以及 domain_strategy 为 ipv6_only 的节点不会分配给该平台。" },
      { type: "p", text: "节点名规则按行生效：普通行满足其一即可，* 开头必须匹配，! 开头表示排除。" },
      { type: "p", text: "租约上的节点失效后，优先换成同一国家的可路由节点。没有同国家节点时，再按平台的分配策略选择。" },
    ],
  },
  {
    id: "circuit",
    title: "熔断",
    blocks: [
      { type: "p", text: "连通性正常但业务不可用时，可以立即熔断节点。节点详情里有「打开熔断」，也可以调用接口。" },
      { type: "code", text: "POST /api/v1/nodes/<节点哈希>/actions/open-circuit\nPOST /<代理令牌>/api/v1/nodes/<节点哈希>/actions/open-circuit" },
    ],
  },
  {
    id: "logs",
    title: "请求日志与登录",
    blocks: [
      { type: "p", text: "请求日志页的「清除日志」会丢掉尚未写入的队列，并删除已经保存的全部请求日志。此操作不能撤销。" },
      { type: "p", text: "管理员令牌失效后，管理接口返回 401，页面会清除本地令牌并回到登录页，同时保留原来的路径。" },
    ],
  },
];

const en: DocSection[] = [
  {
    id: "version",
    title: "Version",
    blocks: [
      { type: "p", text: "The sidebar version is compiled into the binary. docker compose up -d --build reads the VERSION file when RESIN_VERSION is unset. The current value is v1.3.0." },
      { type: "p", text: "Set RESIN_VERSION to override a build. Release builds still use the Git tag." },
    ],
  },
  {
    id: "subscriptions",
    title: "Subscriptions",
    blocks: [
      { type: "p", text: "A remote subscription can set its own User-Agent. Blank uses clash.meta. Local subscriptions cannot set a User-Agent." },
      { type: "p", text: "A blank probe interval uses the global setting. A custom interval applies to that subscription's nodes and must be at least 10 seconds. A node in several subscriptions uses the shortest positive interval." },
      { type: "p", text: "Batch import reads one URL per line. Blank lines and lines starting with # are skipped. Names use the middle domain label unless a regex capture group provides one." },
      { type: "ul", items: ["Plain proxy lines accept IP:PORT and IP:PORT:USER:PASS.", "vmess, vless, trojan, ss, http, socks, sing-box, and Clash subscriptions are also supported."] },
    ],
  },
  {
    id: "healthy",
    title: "Healthy node subscription",
    blocks: [
      { type: "p", text: "Export nodes that are currently healthy and have an egress IP. The admin API uses the admin token. Phone clients use the proxy token." },
      { type: "code", text: "GET /<proxy-token>/api/v1/healthy-subscription\nGET /<proxy-token>/api/v1/healthy-subscription?format=sing-box" },
      { type: "p", text: "The default body is a base64 URI list. format=sing-box returns sing-box JSON." },
    ],
  },
  {
    id: "platforms",
    title: "Platforms and routing",
    blocks: [
      { type: "p", text: "IPv4 only keeps IPv6 exits, IPv6 server addresses, and domain_strategy=ipv6_only nodes out of that platform." },
      { type: "p", text: "Node name rules are line oriented: a plain line matches any, a leading * is required, and a leading ! excludes." },
      { type: "p", text: "When a leased node fails, replacement prefers another routable node in the same country. If none exists, the platform allocation policy chooses." },
    ],
  },
  {
    id: "circuit",
    title: "Circuit breaker",
    blocks: [
      { type: "p", text: "Open a circuit immediately when a node connects but cannot serve the business. The node detail page has this action, and so does the API." },
      { type: "code", text: "POST /api/v1/nodes/<hash>/actions/open-circuit\nPOST /<proxy-token>/api/v1/nodes/<hash>/actions/open-circuit" },
    ],
  },
  {
    id: "logs",
    title: "Request logs and login",
    blocks: [
      { type: "p", text: "Clear logs drops queued entries and deletes every stored request log. It cannot be undone." },
      { type: "p", text: "When the admin token is rejected with 401, the UI clears the stored token and returns to login while keeping the previous path." },
    ],
  },
];

export function docsForLocale(isEnglish: boolean): DocSection[] {
  return isEnglish ? en : zh;
}
