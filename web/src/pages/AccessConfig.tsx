import { Navbar } from "../components/Navbar";
import { useParams } from "react-router";
import { OctagonX } from "lucide-react";
import { SelectListContainer } from "../components/SelectListContainer";

export function AccessConfigPage() {

    const {note} = useParams();

    return (
        <>
            <Navbar/>
            <div className="d-flex flex-column align-items-center container-lg">            
                <h1 className="my-2">Access for Note ID: {note}</h1>
                <div className="my-3 d-flex justify-items-center">
                    <button type="button" className="mx-1 btn btn-primary btn-sm">Users</button>
                    <button type="button" className="mx-1 btn btn-primary btn-sm">User Groups</button>
                    <button type="button" className="mx-1 d-flex align-items-center gap-1 btn btn-primary btn-sm">Access Mode: Blacklist <OctagonX size={16}/> </button>
                </div>
                <input className="form-control mr-sm-2 w-75" type="text" placeholder="Search"></input>
                <SelectListContainer/>
            </div>
        </>
    )
}