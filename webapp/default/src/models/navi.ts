


export type ModelNavi = {
    title?: string
    pages: Array<ModelPage>
};

  
export type ModelPage = {
    level?: number
    type: number
    title: string
    link?: string
    parent?: string
    pages: Array<ModelPage>

};