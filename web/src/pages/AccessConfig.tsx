import { Navbar } from "../components/Navbar";
import { useParams } from "react-router";
import { Check, OctagonX } from "lucide-react";
import { SelectListContainer } from "../components/SelectListContainer";
import type { SelectBlockItem } from "../components/SelectListContainer";
import { useState, useEffect } from "react"

import {
    useMutation,
} from '@tanstack/react-query'

// async function fetchUserInfo() {
//   const res = await fetch('/csh/api/v1/users/info');
//   if (!res.ok) throw new Error('Network response error');
//   return res.json();
// }

type Group = {
    group_name: string;
    uuid: string;
    selected: boolean;
}

type User = {
    username: string;
    uuid: string;
    selected: boolean;
}
type GroupList = {
    groups: Group[];
    cursor: string;
}

type UserList = {
    users: User[];
    cursor: string;
}

const GROUPS_ENDPOINT = "/csh/api/v1/groups";
const USERS_ENDPOINT = "/csh/api/v1/users";
const ACCESS_ENDPOINT = "/csh/api/v1/access";

export function AccessConfigPage() {
    const {note} = useParams();
    const [selected, setSelected] = useState<Set<string>>(new Set());
    const [isBlacklist, setBlacklist] = useState<boolean>(true);
    const [isQueryingUsers, setQueryingUsers] = useState<boolean>(true);

    const [items, setItems] = useState<SelectBlockItem[]>([])
    const [cursor, setCursor] = useState<string>("")
    const [search, setSearch] = useState<string>("")

    const fetchGroups = useMutation<GroupList, Error, {prefix: string, cursor?: string}>({
		mutationKey: [`fetch-new-groups`, cursor],
		mutationFn: (payload: any) =>
        fetch(`${GROUPS_ENDPOINT}`, {
            method: "QUERY",
            body: JSON.stringify(payload),
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
        }).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
		onError: (err) => {
			console.log(err)			
		},
        onSuccess: (data) => {
            setItems((prev) => {
                const existing = new Set(prev.map((item) => item.UUID))

                const newItems = data.groups
                .filter((group) => !existing.has(group.uuid))            
                .map((group): SelectBlockItem => ({
                    Name: group.group_name,
                    Type: "groups",
                    UUID: group.uuid,
                    Icon: false,
                    Selected: group.selected,
                }))

                return [...prev, ...newItems]
            })
            setCursor(data.cursor)
        },
	})

    const fetchUsers = useMutation<UserList, Error, {prefix: string, cursor?: string}>({
		mutationKey: [`fetch-new-users`, cursor],
		mutationFn: (payload: any) =>
        fetch(`${USERS_ENDPOINT}`, {
            method: "QUERY",
            body: JSON.stringify(payload),
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
        }).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
		onError: (err) => {
			console.log(err)			
		},
        onSuccess: (data) => {
            setItems((prev) => {
                const existing = new Set(prev.map((item) => item.UUID))

                const newItems = data.users
                .filter((user) => !existing.has(user.uuid))            
                .map((user): SelectBlockItem => ({
                    Name: user.username,
                    Type: "users",
                    UUID: user.uuid,
                    Icon: true,
                    Selected: user.selected,
                }))

                return [...prev, ...newItems]
            })
            setCursor(data.cursor)
        },
	})

    const updateGroup = useMutation<null, Error, {note_id: string, identifier_uuid: string, member_type: string}>({
		mutationKey: [`update-group-access`],
		mutationFn: (payload: any) =>
        fetch(`${ACCESS_ENDPOINT}`, {
            method: "POST",
            body: JSON.stringify(payload),
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
        }).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
		onError: (err) => {
			console.log(err)			
		},
        onSuccess: () => {
			console.log("Successfully updated group")			
        },
	})

    const toggle = (uuid: string, type: string) => {
        setSelected((prev) => {
            const next = new Set(prev);
            next.has(uuid) ? next.delete(uuid) : next.add(uuid);
            console.log(uuid)

            if (type === "groups") {
                updateGroup.mutate({note_id: "test", identifier_uuid: uuid, member_type: type})
            }
            return next;
        });
    };


    const toggleAccessMode = () => {
        setBlacklist(!isBlacklist)
    }

    const updateSearchMode = (isUsers: boolean) => {
        if (isUsers === isQueryingUsers) return

        setItems([]);
        setCursor("");
        setQueryingUsers(isUsers)
    }

    const onSearchChange = (value: string) => {        
        setSearch(value)
    }

    useEffect(() => {
        const timer = setTimeout(() => {
            setItems([])
            setCursor("")
            if (isQueryingUsers) {
                fetchUsers.mutate({prefix: search})
            } else {
                fetchGroups.mutate({prefix: search})
            }        
        }, 300)
        return () => clearTimeout(timer)
    }, [isQueryingUsers, search])

    return (
        <>
            <Navbar/>
            <div className="d-flex flex-column align-items-center container-lg">            
                <h1 className="my-2">Access for Note ID: {note}</h1>
                <div className="my-3 d-flex justify-items-center">
                    <button type="button" className="mx-1 btn btn-primary btn-sm d-inline-flex align-items-center gap-1" onClick={() => {updateSearchMode(true)}}>Users{isQueryingUsers && <Check size={16}/>}</button>
                    <button type="button" className="mx-1 btn btn-primary btn-sm d-inline-flex align-items-center gap-1" onClick={() => {updateSearchMode(false)}}>User Groups{!isQueryingUsers && <Check size={16}/>}</button>
                    <button type="button" className="mx-1 d-flex align-items-center gap-1 btn btn-primary btn-sm" onClick={toggleAccessMode}>Access Mode: {isBlacklist ? <>Blacklist<OctagonX size={16}/></> : <>Whitelist<Check size={16}/></>}  </button>
                </div>
                <input onChange={(e) => onSearchChange(e.target.value)} className="form-control mr-sm-2 w-75" type="text" placeholder="Search"></input>
                <SelectListContainer items={items} invertSelection={isBlacklist} selected={selected} onToggle={toggle}/>
                <h3><i>Loading...</i></h3>
            </div>
        </>
    )
}