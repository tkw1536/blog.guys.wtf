import Link from "next/link";
import * as React from "react";
import DateDate from "../lib/date";
import { Post } from "../lib/posts";
import Meta from "./meta";

export default class Posting extends React.Component<{ post: Post }> {
    render() {
        const { post: { title = "Untitled Post", date, author, html } } = this.props;
        const dateDate = date ? DateDate.parse(date) : undefined;
        return <article itemScope itemType="https://schema.org/Blog">
            <Meta title={`${title} | ${process.env.NEXT_PUBLIC_BLOG_TITLE}`} />
            <header>
                <h1 itemProp="name headline">{title}</h1>
                {date && <span>
                    <time dateTime={dateDate.toISOString()} itemProp="datePublished">{dateDate.toString()}</time>
                </span>}
                {author && <span itemProp="author" itemScope itemType="https://schema.org/Person">
                    <span itemProp="name">{author}</span>
                </span>}
                
                <nav aria-label="Navigation" itemScope itemType="https://schema.org/SiteNavigationElement">
                    <span><Link href="/" itemProp="url">Home</Link></span>
                    <span><Link legacyBehavior href="/about/" itemProp="url">About</Link></span>
                </nav>
            </header>
            <main dangerouslySetInnerHTML={{ __html: html }} />
        </article>;
    }
}