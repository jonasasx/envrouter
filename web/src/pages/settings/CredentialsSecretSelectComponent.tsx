import {Box, Button, CircularProgress, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import * as React from "react";
import {useEffect, useState} from "react";
import {OptionsObject, useSnackbar} from "notistack";
import NewCredentialsSecretComponent from "./NewCredentialsSecretComponent";
import {CredentialsSecretListItem, DefaultApiFp} from "../../axios";

interface CredentialsSecretSelectProps {
    secret: string
    onSecretChange: (secret: string) => void
}

const api = DefaultApiFp()
const loadAllSecrets = (
    setLoading: ((b: boolean) => void),
    setCredentialsSecrets: ((r: Array<CredentialsSecretListItem>) => void),
    enqueueSnackbar: ((s: string, options?: OptionsObject) => void)
) => {
    setLoading(true)
    api.apiV1CredentialsSecretsGet()
        .then(request => request())
        .then(response => {
            setCredentialsSecrets(response.data)
        })
        .then(() => setLoading(false))
        .catch(() => {
            setLoading(false)
            enqueueSnackbar(`Loading failed`, {variant: "error"})
        })
}

export default function CredentialsSecretSelectComponent(props: CredentialsSecretSelectProps) {
    const [credentialsSecrets, setCredentialsSecrets] = useState<Array<CredentialsSecretListItem>>([])
    const [loading, setLoading] = useState(false)
    const [openCreateNewSecret, setOpenCreateNewSecret] = useState(false)
    const [selectedSecret, setSelectedSecret] = useState<string>(props.secret || "")
    const {onSecretChange} = props
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();


    useEffect(() => {
        loadAllSecrets(setLoading, setCredentialsSecrets, enqueueSnackbar)
        return () => {
            closeSnackbar()
        }
    }, [enqueueSnackbar, closeSnackbar])
    return (
        <React.Fragment>
            <Box
                style={{outline: "none"}}
                component="form"
                noValidate
                autoComplete="off"
            >
                {loading && <CircularProgress/> ||
                    <FormControl variant="standard" fullWidth sx={{mb: 1, mt: 1}}>
                        <InputLabel id="secret-select-standard-label">Repository credentials Secret</InputLabel>
                        <Select
                            labelId="secret-select-standard-label"
                            id="secret-select-standard"
                            value={selectedSecret}
                            onChange={(e) => {
                                setSelectedSecret(e.target.value)
                                onSecretChange(e.target.value)
                            }}
                            label="Repository credentials Secret"
                        >
                            <MenuItem value="">
                                <em>None</em>
                            </MenuItem>
                            {credentialsSecrets.map(credentialsSecret => (
                                <MenuItem key={credentialsSecret.name} value={credentialsSecret.name}>{credentialsSecret.name}</MenuItem>
                            ))}
                        </Select>
                    </FormControl>}
            </Box>
            <p>Or press the button to create:</p>

            <Button variant="outlined" onClick={() => {
                setOpenCreateNewSecret(true)
            }}>Create new credentials secret</Button>

            {openCreateNewSecret && <NewCredentialsSecretComponent onClose={(credentialsSecret) => {
                setOpenCreateNewSecret(false)
                if (credentialsSecret !== undefined) {
                    setCredentialsSecrets([...credentialsSecrets, credentialsSecret])
                    setSelectedSecret(credentialsSecret.name)
                    onSecretChange(credentialsSecret.name)
                }
            }}/>}
        </React.Fragment>
    )
}