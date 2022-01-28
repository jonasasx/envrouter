import {useState} from "react";
import {Repository} from "../../axios";
import RepositoryGridComponent from "./RepositoryGridComponent";

export default function SettingsPage() {
    const [repositories, setRepositories] = useState<Array<Repository>>([])
    return (
        <div>
            <RepositoryGridComponent onRepositoriesUpdate={repositories => setRepositories(repositories)}/>
        </div>
    )
}