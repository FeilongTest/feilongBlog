interface ArticleList {
    ID:number;
    content:string;
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
}
interface ArticelSummary {
    fid:number;
    total:number;
    name:string;
}
export type {
    ArticleList,
    ArticelSummary,
}