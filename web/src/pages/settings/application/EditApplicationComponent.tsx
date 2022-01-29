import {WithStyles, withStyles} from "@mui/styles";
import {Application, DefaultApiFp, Repository} from "../../../axios";
import {Theme} from "@mui/material/styles";
import {useEffect, useState} from "react";
import {Box, Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, Select, TextField} from "@mui/material";
import {useSnackbar} from "notistack";

interface EditApplicationProps extends WithStyles<typeof styles> {
    application: Application
    repositories: Array<Repository>
    onClose: (app: Application | undefined) => void
}

const styles = (theme: Theme) => ({})
const api = DefaultApiFp()
export default withStyles(styles)(function EditApplicationComponent(props: EditApplicationProps) {
    const {classes, application, repositories, onClose} = props
    const [loading, setLoading] = useState(false)
    const name = application.name
    const [repositoryName, setRepositoryName] = useState(application.repositoryName)
    const [webhook, setWebhook] = useState(application.webhook)
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();
    const save = () => {
        const newApplication = {name, repositoryName, webhook} as Application
        api.apiV1ApplicationsNamePut(application.name, newApplication)
            .then(request => request())
            .then(response => {
                onClose(response.data)
                setLoading(false);
            })
            .catch(() => {
                setLoading(false)
                enqueueSnackbar(`Saving failed`, {variant: "error"})
            })
    }
    useEffect(() => {
        return () => {
            closeSnackbar()
        }
    }, [closeSnackbar])
    return (
        <Dialog
            open={true}
            onClose={() => !loading && onClose(undefined)}
            fullWidth={true}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
        >
            <DialogTitle id="alert-dialog-title">
                Configure {name}
            </DialogTitle>
            <DialogContent>
                <FormControl variant="standard" fullWidth>
                    <InputLabel id="repo-select-standard-label">Repository</InputLabel>
                    <Select
                        labelId="repo-select-standard-label"
                        id="repo-select-standard"
                        value={repositoryName}
                        onChange={(e) => {
                            setRepositoryName(e.target.value)
                        }}
                        label="Repository"
                    >
                        <MenuItem value="">
                            <em>None</em>
                        </MenuItem>
                        {repositories.map(repository => (
                            <MenuItem value={repository.name}>{repository.name}</MenuItem>
                        ))}
                    </Select>
                </FormControl>
                <TextField label="Webhook URL" variant="standard" value={webhook} onChange={e => setWebhook(e.target.value)} fullWidth/>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => !loading && onClose(undefined)} disabled={loading}>Cancel</Button>
                <Box sx={{m: 1, position: 'relative'}}>
                    <Button onClick={() => {
                        setLoading(true)
                        save()
                    }} autoFocus disabled={loading}>Save</Button>
                    {loading && (
                        <CircularProgress
                            size={24}
                            sx={{
                                position: 'absolute',
                                top: '50%',
                                left: '50%',
                                marginTop: '-12px',
                                marginLeft: '-12px',
                            }}/>)}
                </Box>
            </DialogActions>
        </Dialog>
    )
})