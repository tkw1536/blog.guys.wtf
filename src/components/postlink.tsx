import * as React from 'react';
import Link from "next/link";
import { Post } from "../lib/posts";
import DateDate from "../lib/date";

export default class PostLink extends  React.Component<{ post: Post}> {
    render() {
        const { post: { path, title = "Untitled Post", date, author} } = this.props;
        return <>
            <Link href={`/${path}`}>{title}</Link>
            {date && <span>{DateDate.parse(date).toString()}</span>}
            {author && <span>{author}</span>}
        </>
    }
}