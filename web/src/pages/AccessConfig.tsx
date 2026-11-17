import { Navbar } from "../components/Navbar";
import { useParams } from "react-router";
import { OctagonX } from "lucide-react";
import { SelectListContainer } from "../components/SelectListContainer";
import { useState } from "react"


// async function fetchUserInfo() {
//   const res = await fetch('/csh/api/v1/users/info');
//   if (!res.ok) throw new Error('Network response error');
//   return res.json();
// }

export function AccessConfigPage() {
    const {note} = useParams();
    const [selected, setSelected] = useState<Set<string>>(new Set());
    const [isBlacklist, setBlacklist] = useState<boolean>(true);


    const toggle = (uuid: string) => {
        setSelected((prev) => {
            const next = new Set(prev);
            next.has(uuid) ? next.delete(uuid) : next.add(uuid);
            return next;
        });
    };


    const toggleAccessMode = () => {
        setBlacklist(!isBlacklist)
    }

    return (
        <>
            <Navbar/>
            <div className="d-flex flex-column align-items-center container-lg">            
                <h1 className="my-2">Access for Note ID: {note}</h1>
                <div className="my-3 d-flex justify-items-center">
                    <button type="button" className="mx-1 btn btn-primary btn-sm">Users</button>
                    <button type="button" className="mx-1 btn btn-primary btn-sm">User Groups</button>
                    <button type="button" className="mx-1 d-flex align-items-center gap-1 btn btn-primary btn-sm" onClick={toggleAccessMode}>Access Mode: {isBlacklist ? <>Blacklist<OctagonX size={16}/></> : <>Whitelist<OctagonX size={16}/></>}  </button>
                </div>
                <input className="form-control mr-sm-2 w-75" type="text" placeholder="Search"></input>
                <SelectListContainer invertSelection={isBlacklist} selected={selected} onToggle={toggle}/>
            </div>
        </>
    )
}