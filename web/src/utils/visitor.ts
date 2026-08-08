const VISITOR_KEY = "blog_visitor_id";
const LIKED_KEY = "blog_liked_articles";

export const getVisitorId = () => {
  let visitorId = localStorage.getItem(VISITOR_KEY);
  if (!visitorId) {
    visitorId = crypto.randomUUID();
    localStorage.setItem(VISITOR_KEY, visitorId);
  }
  return visitorId;
};

export const getLikedArticles = () => {
  try {
    return new Set<number>(JSON.parse(localStorage.getItem(LIKED_KEY) || "[]"));
  } catch {
    return new Set<number>();
  }
};

export const saveArticleLike = (articleId: number, liked: boolean) => {
  const articles = getLikedArticles();
  liked ? articles.add(articleId) : articles.delete(articleId);
  localStorage.setItem(LIKED_KEY, JSON.stringify([...articles]));
};
