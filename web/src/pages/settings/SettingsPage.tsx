import {useState} from "react";
import {Repository} from "../../axios";
import RepositoryGridComponent from "./repository/RepositoryGridComponent";
import ApplicationGridComponent from "./application/ApplicationGridComponent";

export default function SettingsPage() {
    const [repositories, setRepositories] = useState<Array<Repository>>([])
    return (
        <div>
            <RepositoryGridComponent onRepositoriesUpdate={repositories => setRepositories(repositories)}/>
            <p>&nbsp;</p>
            <ApplicationGridComponent repositories={repositories}/>
        </div>
    )
}