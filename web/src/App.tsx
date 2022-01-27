import React from 'react';
import './App.css';
import {SnackbarProvider} from 'notistack';
import { ThemeProvider} from '@mui/styles';
import DashboardPage from './pages/dashboard/DashboardPage';
import {NavLink, Route, Routes, useLocation} from 'react-router-dom';
import SettingsPage from './pages/settings/SettingsPage';
import {Box, createTheme, CssBaseline, IconButton, ListItem, MenuItem, MenuList, Typography} from '@mui/material';
import {CSSObject, styled, Theme} from '@mui/material/styles';
import MuiAppBar, {AppBarProps as MuiAppBarProps} from '@mui/material/AppBar';
import MuiDrawer from '@mui/material/Drawer';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import Toolbar from '@mui/material/Toolbar';
import Divider from '@mui/material/Divider';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import MenuIcon from '@mui/icons-material/Menu';
import DashboardIcon from '@mui/icons-material/Dashboard';
import SourceIcon from '@mui/icons-material/Source';

const theme = createTheme({
  palette: {
    background: {
      default: "#f5f5f5"
    }
  }
});

const drawerWidth = 240;

const openedMixin = (theme: Theme): CSSObject => ({
  width: drawerWidth,
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
  overflowX: 'hidden',
});

const closedMixin = (theme: Theme): CSSObject => ({
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  overflowX: 'hidden',
  width: `calc(${theme.spacing(7)} + 1px)`,
  [theme.breakpoints.up('sm')]: {
    width: `calc(${theme.spacing(9)} + 1px)`,
  },
});

const DrawerHeader = styled('div')(({theme}) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
}));

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({theme, open}) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, {shouldForwardProp: (prop) => prop !== 'open'})(
    ({theme, open}) => ({
      width: drawerWidth,
      flexShrink: 0,
      whiteSpace: 'nowrap',
      boxSizing: 'border-box',
      ...(open && {
        ...openedMixin(theme),
        '& .MuiDrawer-paper': openedMixin(theme),
      }),
      ...(!open && {
        ...closedMixin(theme),
        '& .MuiDrawer-paper': closedMixin(theme),
      }),
    }),
);

function App() {
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };
  const location = useLocation();

  const activeRoute = (routeName: any) => {
    return location.pathname === routeName ? true : false;
  }
  return (
      <div className="App">
        <ThemeProvider theme={theme}>
          <SnackbarProvider maxSnack={3} anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}>

            <Box sx={{display: 'flex'}}>
              <CssBaseline/>
              <AppBar position="fixed" open={open}>
                <Toolbar>
                  <IconButton
                      color="inherit"
                      aria-label="open drawer"
                      onClick={handleDrawerOpen}
                      edge="start"
                      sx={{
                        marginRight: '36px',
                        ...(open && {display: 'none'}),
                      }}
                  >
                    <MenuIcon/>
                  </IconButton>
                  <Typography variant="h6" noWrap component="div">
                    Env Router
                  </Typography>
                </Toolbar>
              </AppBar>
              <Drawer variant="permanent" open={open}>
                <DrawerHeader>
                  <IconButton onClick={handleDrawerClose}>
                    {theme.direction === 'rtl' ? <ChevronRightIcon/> : <ChevronLeftIcon/>}
                  </IconButton>
                </DrawerHeader>
                <Divider/>
                <MenuList>
                  {[
                    {title: 'Dashboard', icon: <DashboardIcon/>, path: '/'},
                    {title: 'Repositories', icon: <SourceIcon/>, path: '/repo'}
                  ].map((item, index) => (
                      <NavLink to={item.path} style={{textDecoration: 'none'}} key={item.path}>
                        <ListItem key={item.title} selected={activeRoute(item.path)}>
                          <ListItemIcon>
                            {item.icon}
                          </ListItemIcon>
                          <ListItemText primary={item.title}/>
                        </ListItem>
                      </NavLink>
                  ))}
                </MenuList>
                <Divider/>
              </Drawer>
              <Box component="main" sx={{flexGrow: 1, p: 3}}>
                <DrawerHeader/>

                <Routes>
                  <Route path="/repo" element={<SettingsPage/>}/>
                  <Route path="/" element={<DashboardPage/>}/>
                </Routes>
              </Box>
            </Box>
          </SnackbarProvider>
        </ThemeProvider>
      </div>
  );
}

export default App;
