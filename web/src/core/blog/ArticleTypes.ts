interface ArticleList {
    ID:number;
    content:string;
    contentFormat:"html" | "markdown";
    fid:number;
    pic:string;
    title:string;
    uid:number;
    view:number;
    type:number;
    ctime:number;
    editTime:number;
    file:string;
    isTop:number;
    status:number;
    commentCount:number;
    likeCount:number;
}
interface ArticelSummary {
    fid:number;
    total:number;
    name:string;
}

interface ArticleBrief {
    ID:number;
    fid:number;
    title:string;
    pic:string;
    view:number;
    ctime:number;
}

interface PostAside {
    related:ArticleBrief[];
    hot:ArticleBrief[];
}

export type {
    ArticleList,
    ArticelSummary,
    ArticleBrief,
    PostAside,
}
