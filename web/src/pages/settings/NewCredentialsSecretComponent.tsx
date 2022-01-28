import {Box, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";
import * as React from "react";
import {useState} from "react";
import Button from "@mui/material/Button";
import {useSnackbar} from "notistack";
import {CredentialsSecretListItem, CredentialsSecretRequest, DefaultApiFp} from "../../axios";

interface NewCredentialsSecretProps {
    onClose: (secret: CredentialsSecretListItem | undefined) => void
}

const api = DefaultApiFp()
export default function NewCredentialsSecretComponent(props: NewCredentialsSecretProps) {
    const [loading, setLoading] = useState(false)
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [key, setKey] = useState("")
    const {onClose} = props
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();

    const save = () => {
        setLoading(true)
        const credentialsSecretRequest = {username, password, key} as CredentialsSecretRequest
        api.apiV1CredentialsSecretsPost(credentialsSecretRequest)
            .then(request => request())
            .then(response => response.data)
            .then((credentialsSecret) => {
                setLoading(false)
                onClose(credentialsSecret)
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
                Create new credentials secret
            </DialogTitle>
            <DialogContent>
                <Box
                    component="form"
                    sx={{
                        '& .MuiTextField-root': {mb: 1},
                    }}
                    noValidate
                    autoComplete="off"
                >
                    <TextField variant="standard" fullWidth label="Username" value={username} onChange={(e) => setUsername(e.target.value)}/>
                    <TextField variant="standard" fullWidth label="Password" value={password} onChange={(e) => setPassword(e.target.value)}/>
                    <TextField variant="standard" fullWidth multiline rows={10} label="TLS private key" value={key} onChange={(e) => setKey(e.target.value)}/>
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => !loading && onClose(undefined)} disabled={loading}>Cancel</Button>
                <Button onClick={() => !loading && save()} disabled={loading}>Save</Button>
            </DialogActions>
        </Dialog>
    )
}