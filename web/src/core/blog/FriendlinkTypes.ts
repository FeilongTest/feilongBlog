interface Friendlink {
    ID: number;
    name: string;
    url: string;
    logo: string;
    desc: string;
    status: number;
    ctime: number;
}

interface FriendlinkFormData {
    ID?: number;
    name: string;
    url: string;
    logo: string;
    desc: string;
    status: number;
}

export type {
    Friendlink,
    FriendlinkFormData,
}

