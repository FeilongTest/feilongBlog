export const delHtmlTag = (str) => {
    return str.replace(/<[^>]+>/g, '').replaceAll('&amp;','&').replaceAll('&nbsp;',' ')
}