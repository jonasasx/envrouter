import {Box, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";
import Button from "@mui/material/Button";
import * as React from "react";
import {useState} from "react";
import {useSnackbar} from "notistack";
import {DefaultApiFp, Repository} from "../../../axios";
import {withStyles, WithStyles} from "@mui/styles";
import {Theme} from "@mui/material/styles";
import CredentialsSecretSelectComponent from "./CredentialsSecretSelectComponent";

interface EditRepositoryProps extends WithStyles<typeof styles> {
    repository: Repository
    onClose: (repo: Repository | undefined) => void
}

const styles = (theme: Theme) => ({})

const api = DefaultApiFp()

export default withStyles(styles)(function EditRepositoryComponent(props: EditRepositoryProps) {
    const {repository, onClose} = props
    const [loading, setLoading] = useState(false)
    const [name, setName] = useState(repository.name || '')
    const [url, setUrl] = useState(repository.url || '')
    const [credentialsSecret, setCredentialsSecret] = useState(repository.credentialsSecret)
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();
    const save = () => {
        setLoading(true)
        const newRepository = {name, url, credentialsSecret} as Repository
        api.apiV1RepositoriesPost(newRepository)
            .then(request => request())
            .then(response => response.data)
            .then(savedRepository => {
                setLoading(false);
                onClose(savedRepository)
            })
            .catch(() => {
                setLoading(false)
                enqueueSnackbar(`Saving failed`, {variant: "error"})
            })
    }
    return (
        <Dialog
            open={true}
            onClose={() => !loading && onClose(undefined)}
            fullWidth={true}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
        >
            <DialogTitle id="alert-dialog-title">
                Configure repository {repository.isNew && '' || name}
            </DialogTitle>
            <DialogContent>
                {repository.isNew && <TextField label="git repo name" variant="standard" value={name} onChange={e => setName(e.target.value)} fullWidth/>}
                <TextField label="git repo url" variant="standard" value={url} onChange={e => setUrl(e.target.value)} fullWidth/>
                <CredentialsSecretSelectComponent secret={credentialsSecret} onSecretChange={(credentialsSecret) => setCredentialsSecret(credentialsSecret)}/>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => !loading && onClose(undefined)} disabled={loading}>Cancel</Button>
                <Box sx={{m: 1, position: 'relative'}}>
                    <Button onClick={() => {
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