### PtHackPoc

相关博客: https://blog.lyc8503.site/post/pt-hack/

实现功能:

1. 免费下载种子
   
   用法: 直接无需参数运行主程序, 首次运行时会自动读取当前目录下所有 torrent 文件, 删除原种子并生成 `FREE_` 打头的新种子, 此时在同一台电脑上再启动 BT 客户端导入这些新种子就能免费下载.

2. 刷 peer 流量

   用法: 运行主程序时带 peer IP:port 和 info_hash 两个参数, 例如 `./ptcheat '[2001:da8:1007:4000::1:4eff]:33006' ecfdd75b8b493c6f0cb7142bba66466183f0f707`

(特别是在错误使用的情况下)有封号风险, 后果自负.

---

另: py 文件夹下是我写的一个 PT 站 RSS 爬虫, 此处附 BYR 和 tjupt 的爬虫结果, 说不定有用. (比如在不访问 PT 站的情况下查 info_hash)

https://pan.lyc8503.site/Public/%E5%85%B6%E4%BB%96/PT