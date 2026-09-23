import { Link } from "react-router-dom";
import guideEn from "./user-guide.en.md?raw";
import guideZh from "./user-guide.zh-CN.md?raw";
import { LanguageSwitcher } from "../../components/LanguageSwitcher";
import { useAuthStore } from "../auth/auth-store";
import { useI18n } from "../../i18n";
import { renderMarkdown } from "./renderMarkdown";

export function DocsPage() {
  const { isEnglish, t } = useI18n();
  const token = useAuthStore((state) => state.token);
  const markdown = isEnglish ? guideEn : guideZh;

  return (
    <main className="docs-public">
      <header className="docs-public-bar">
        <div>
          <p className="docs-public-kicker">Resin</p>
          <h1>{t("使用文档")}</h1>
        </div>
        <div className="docs-public-actions">
          <a href={isEnglish ? "/ui/user-guide.en.md" : "/ui/user-guide.zh-CN.md"}>{t("Markdown")}</a>
          <a href="/ui/llms.txt">llms.txt</a>
          <LanguageSwitcher />
          <Link to={token ? "/dashboard" : "/login"}>{token ? t("进入控制台") : t("管理员登录")}</Link>
        </div>
      </header>
      <article className="docs-article docs-public-article">{renderMarkdown(markdown)}</article>
    </main>
  );
}
