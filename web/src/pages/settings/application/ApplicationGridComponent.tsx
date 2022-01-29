import {WithStyles, withStyles} from "@mui/styles";
import EditIcon from '@mui/icons-material/Edit';
import {Application, DefaultApiFp, Repository} from "../../../axios";
import {Theme} from "@mui/material/styles";
import {DataGrid, GridActionsCellItem, GridColumns,} from '@mui/x-data-grid';
import {Container, Paper} from "@mui/material";
import {useEffect, useState} from "react";
import {OptionsObject, useSnackbar} from "notistack";
import EditApplicationComponent from "./EditApplicationComponent";

interface ApplicationGridProps extends WithStyles<typeof styles> {
    repositories: Array<Repository>
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
const comparator = (a: Application, b: Application) => a.name.localeCompare(b.name)
const loadAllApplications = (
    setLoading: ((b: boolean) => void),
    setApplications: ((r: Array<Application>) => void),
    enqueueSnackbar: ((s: string, options?: OptionsObject) => void)
) => {
    setLoading(true)
    api.apiV1ApplicationsGet()
        .then(request => request())
        .then(response => setApplications(response.data))
        .then(() => setLoading(false))
        .catch(() => {
            setLoading(false)
            enqueueSnackbar(`Loading failed`, {variant: "error"})
        })
}

export default withStyles(styles)(function ApplicationGridComponent(props: ApplicationGridProps) {
    const {classes, repositories} = props
    const [editingApplication, setEditingApplication] = useState<Application | undefined>(undefined)
    const [applications, setApplications] = useState<Array<Application>>([])
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();
    const [loading, setLoading] = useState(false)

    const updateApplication = (application: Application) => {
        setApplications(currentApplications => {
            const index = currentApplications.findIndex(a => a.name === application.name)
            const result = (index === -1) &&
                [...currentApplications, application] ||
                [
                    ...currentApplications.slice(0, index),
                    application,
                    ...currentApplications.slice(index + 1)
                ]
            return result.sort(comparator)
        })
    }

    useEffect(() => {
        loadAllApplications(setLoading, setApplications, enqueueSnackbar)

        return () => {
            closeSnackbar()
        }
    }, [enqueueSnackbar, closeSnackbar])

    const columns: GridColumns = [
        {field: 'name', headerName: 'Application Name', flex: 1, type: 'string', editable: false},
        {field: 'repositoryName', headerName: 'Repository Name', flex: 1, type: 'string', editable: false},
        {field: 'webhook', headerName: 'Webhook', flex: 1, type: 'string', editable: false},
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
                            setEditingApplication(row as Application)
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
                    rows={applications}
                    loading={loading}
                    columns={columns}
                    components={{}}
                />
                {editingApplication && <EditApplicationComponent application={editingApplication} repositories={repositories} onClose={(application: Application | undefined) => {
                    setEditingApplication(undefined)
                    application && updateApplication(application)
                }}/>}
            </Paper>
        </Container>
    );
})