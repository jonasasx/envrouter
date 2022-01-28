import {Button, Container, Paper} from "@mui/material";
import {DataGrid, GridActionsCellItem, GridColumns, GridToolbarContainer,} from '@mui/x-data-grid';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/DeleteOutlined';
import {useEffect, useState} from "react";
import {DefaultApiFp, Repository} from "../../axios";
import {OptionsObject, useSnackbar} from "notistack";
import {WithStyles, withStyles} from "@mui/styles";
import {Theme} from "@mui/material/styles";
import EditRepositoryComponent from "./EditRepositoryComponent";

interface RepositoryGrid extends WithStyles<typeof styles> {
    onRepositoriesUpdate: (repositories: Array<Repository>) => void
}

const styles = (theme: Theme) => ({
    actions: {
        color: theme.palette.text.secondary,
    },
    textPrimary: {
        color: theme.palette.text.primary,
    }
})

const api = DefaultApiFp()
const comparator = (a: Repository, b: Repository) => a.name.localeCompare(b.name)
const loadAllRepositories = (
    setLoading: ((b: boolean) => void),
    setRepositories: ((r: Array<Repository>) => void),
    enqueueSnackbar: ((s: string, options?: OptionsObject) => void)
) => {
    console.log('loadAllRepositories')
    setLoading(true)
    api.apiV1RepositoriesGet()
        .then(request => request())
        .then(response => {
            setRepositories(response.data.sort(comparator))
            setLoading(false)
        })
        .catch(() => {
            setLoading(false)
            enqueueSnackbar(`Loading failed`, {variant: "error"})
        })
}

export default withStyles(styles)(function RepositoryGridComponent(props: RepositoryGrid) {
    const {classes, onRepositoriesUpdate} = props
    const [loading, setLoading] = useState(false)
    const [repositories, setRepositories] = useState<Array<Repository>>([])
    const [editingRepository, setEditingRepository] = useState<Repository | undefined>(undefined)
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();

    const deleteRepository = (repository: Repository) => {
        setLoading(true)
        api.apiV1RepositoriesNameDelete(repository.name)
            .then(request => request())
            .then(() => {
                setRepositories((currentRepositories) => {
                    const index = currentRepositories.findIndex(r => r.name === repository.name)
                    if (index !== -1) {
                        return [
                            ...currentRepositories.slice(0, index),
                            ...currentRepositories.slice(index + 1)
                        ]
                    }
                    return currentRepositories
                })
                setLoading(false)
            })
            .catch(() => {
                setLoading(false)
                enqueueSnackbar(`Loading failed`, {variant: "error"})
            })
    }

    const addRepository = (repository: Repository) => {
        setRepositories(currentRepositories => {
            const index = currentRepositories.findIndex(r => r.name === repository.name)
            const result = (index === -1) &&
                [...currentRepositories, repository] ||
                [
                    ...currentRepositories.slice(0, index),
                    repository,
                    ...currentRepositories.slice(index + 1)
                ]
            return result.sort(comparator)
        })
    }

    useEffect(() => {
        onRepositoriesUpdate(repositories)
    }, [repositories, onRepositoriesUpdate])

    useEffect(() => {
        loadAllRepositories(setLoading, setRepositories, enqueueSnackbar)
        return () => {
            closeSnackbar()
        }
    }, [enqueueSnackbar, closeSnackbar])

    const columns: GridColumns = [
        {field: 'name', headerName: 'Repository Name', flex: 1, type: 'string', editable: false},
        {field: 'url', headerName: 'Repository Url', flex: 1, type: 'string', editable: false},
        {
            field: 'actions',
            type: 'actions',
            headerName: 'Actions',
            editable: false,
            width: 100,
            cellClassName: classes.actions,
            getActions: ({id, row}) => {
                return [
                    <GridActionsCellItem
                        icon={<EditIcon/>}
                        label="Edit"
                        className={classes.textPrimary}
                        onClick={(e) => {
                            setEditingRepository(row as Repository)
                        }}
                        color="inherit"
                    />,
                    <GridActionsCellItem
                        icon={<DeleteIcon/>}
                        label="Delete"
                        onClick={() => {
                            deleteRepository(row as Repository)
                        }}
                        color="inherit"
                    />,
                ];
            },
        },
    ];

    return (
        <Container>
            <Paper>
                <DataGrid
                    getRowId={(repo) => (repo.name)}
                    autoHeight
                    rows={repositories}
                    loading={loading}
                    columns={columns}
                    components={{
                        Toolbar: () => (
                            <GridToolbarContainer>
                                <Button color="primary" startIcon={<AddIcon/>} onClick={() => setEditingRepository({isNew: true} as Repository)}>
                                    Add new Repository
                                </Button>
                            </GridToolbarContainer>
                        ),
                    }}
                />
                {editingRepository && <EditRepositoryComponent repository={editingRepository} onClose={(repository: Repository | undefined) => {
                    setEditingRepository(undefined)
                    repository && addRepository(repository)
                }}/>}
            </Paper>
        </Container>
    )
})