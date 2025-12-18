import menu from "@/router/menu";

/**
 * 根据当前 path 生成面包屑
 * @param {string} path
 * @returns [{ label, path }]
 */
export function getBreadcrumb(path) {
    const result = [];

    function dfs(list, parents = []) {
        for (const item of list) {
            const current = [...parents, item];

            if (item.path === path) {
                result.push(...current);
                return true;
            }

            if (item.children && dfs(item.children, current)) {
                return true;
            }
        }
        return false;
    }

    dfs(menu);
    return result;
}
