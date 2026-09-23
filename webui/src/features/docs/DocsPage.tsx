import { BookOpen } from "lucide-react";
import { docsForLocale } from "./docsContent";
import { useI18n } from "../../i18n";

export function DocsPage() {
  const { isEnglish, t } = useI18n();
  const sections = docsForLocale(isEnglish);

  return (
    <section className="docs-page">
      <header className="page-header">
        <div>
          <h2>{t("使用文档")}</h2>
          <p className="module-description">{t("当前版本的操作说明。切换语言会切换文档语言。")}</p>
        </div>
      </header>

      <div className="docs-layout">
        <nav className="docs-toc" aria-label={t("文档目录")}>
          {sections.map((section) => (
            <a key={section.id} href={`#${section.id}`}>
              {section.title}
            </a>
          ))}
        </nav>

        <article className="docs-article">
          {sections.map((section) => (
            <section key={section.id} id={section.id} className="docs-section">
              <h3>
                <BookOpen size={16} aria-hidden="true" />
                {section.title}
              </h3>
              {section.blocks.map((block, index) => {
                if (block.type === "ul") {
                  return (
                    <ul key={index}>
                      {block.items.map((item) => (
                        <li key={item}>{item}</li>
                      ))}
                    </ul>
                  );
                }
                if (block.type === "code") {
                  return <pre key={index}><code>{block.text}</code></pre>;
                }
                return <p key={index}>{block.text}</p>;
              })}
            </section>
          ))}
        </article>
      </div>
    </section>
  );
}
