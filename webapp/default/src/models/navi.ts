


export type ModelNavi = {
    title?: string
    pages: Array<ModelPage>
    header: Array<ModelPage>
    footer: Array<ModelPage>
    subnavis: Array<ModelNavi>
    id: string
};

export type ModelContent = {
    html:string
    toc: Array<ModelToc>
    id: string
};


  
export type ModelPage = {
    show: boolean

    id: string
    level?: number
    type: number
    title: string
    link?: string
    parent?: string
    pages: Array<ModelPage>


};
export type ModelToc = {
    id: string
    title: string
    items: Array<ModelToc>

};