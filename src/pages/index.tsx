import * as React from 'react';
import { GetStaticProps } from "next";
import Meta from "../components/meta";
import Link from "next/link";
import { Post, readAll } from "../lib/posts";
import PostLink from "../components/postlink";

export default class Page extends React.Component<{ posts: Array<Post> }> {
    render() {
        const { posts } = this.props;
        return <article>
            <Meta title={process.env.NEXT_PUBLIC_BLOG_TITLE} />
            <header>
                <h1>{process.env.NEXT_PUBLIC_BLOG_TITLE}</h1>
                <span>{posts.length} Posts</span>
                <span><Link href="/about/">About</Link></span>
            </header>
            <main>
                <ul className="posts">
                    {posts.map(post => <li key={post.path}><PostLink post={post} /></li>)}
                </ul>
            </main>
        </article>;
    }
}

export const getStaticProps: GetStaticProps = async () => {
    const posts = await readAll();
    return {
        props: {
            posts,
        }
    }
}
