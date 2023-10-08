


export type ModelNavi = {
    title?: string
    pages: Array<ModelPage>
    header: Array<ModelPage>
    footer: Array<ModelPage>
    subnavis: Array<ModelNavi>
    id: string
    accordion: boolean
};

export type ModelContent = {
    html:string
    toc: Array<ModelToc>
    id: string
};


  
export type ModelPage = {
    show: boolean
    activeParent: boolean
    actrive: boolean
    id: string
    level?: number
    type: string
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