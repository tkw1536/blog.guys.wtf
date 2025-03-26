import Link from "next/link";
import * as React from "react";
import Meta from "./meta";

export default class StaticPage extends React.Component<{ title: string, html: string }> {
    render() {
        const { title, html } = this.props;
        return <article itemScope itemType="https://schema.org/WebSite">
            <Meta title={`${title} | ${process.env.NEXT_PUBLIC_BLOG_TITLE}`} />
            <header>
                <h1 itemProp="name headline">{title}</h1>
                <span>{title}</span>
                <nav aria-label="Navigation" itemScope itemType="https://schema.org/SiteNavigationElement">
                    <span><Link href="/" itemProp="url">Home</Link></span>
                </nav>
            </header>
            <main dangerouslySetInnerHTML={{ __html: html }} />
        </article>;
    }
}
