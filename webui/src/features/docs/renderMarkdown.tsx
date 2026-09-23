import type { ReactNode } from "react";

function slug(text: string): string {
  return text
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\u4e00-\u9fff]+/gi, "-")
    .replace(/^-|-$/g, "");
}

export function renderMarkdown(source: string): ReactNode[] {
  const body = source.replace(/^---[\s\S]*?---\s*/, "");
  const blocks = body.split(/\n{2,}/);
  return blocks.map((block, index) => {
    const lines = block.split("\n").filter((line) => line.length > 0);
    if (lines.length === 0) {
      return null;
    }
    if (lines.every((line) => line.startsWith("```"))) {
      return null;
    }
    if (lines[0].startsWith("```")) {
      const code = lines.slice(1, lines.at(-1) === "```" ? -1 : undefined).join("\n");
      return (
        <pre key={index}>
          <code>{code}</code>
        </pre>
      );
    }
    if (lines[0].startsWith("# ")) {
      const text = lines[0].slice(2);
      return <h2 key={index} id={slug(text)}>{text}</h2>;
    }
    if (lines[0].startsWith("## ")) {
      const text = lines[0].slice(3);
      return <h3 key={index} id={slug(text)}>{text}</h3>;
    }
    if (lines[0].startsWith("### ")) {
      const text = lines[0].slice(4);
      return <h4 key={index} id={slug(text)}>{text}</h4>;
    }
    if (lines.every((line) => line.startsWith("- "))) {
      return (
        <ul key={index}>
          {lines.map((line) => (
            <li key={line}>{line.slice(2)}</li>
          ))}
        </ul>
      );
    }
    return <p key={index}>{lines.join(" ")}</p>;
  });
}
