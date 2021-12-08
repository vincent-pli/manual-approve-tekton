import AppBar from '@material-ui/core/AppBar';
import Box from '@material-ui/core/Box';
import IconButton from '@material-ui/core/IconButton';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import DoneOutlineIcon from '@mui/icons-material/DoneOutline';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import * as React from 'react';


class Datatable extends React.Component {
    constructor(props) {
      super(props);
      console.log(props.location.approve)
      this.state = {
        date: new Date(),
        approve: props.location.approve,
      };

    }
  
    componentDidMount() {
      console.log('I was triggered during componentDidMount of Datatable')
    }
  
    approveRequest(row){
        let approveTemplate = this.state.approve.templateNS + "/" + this.state.approve.templateName;
        if(row){
            fetch(window._env_.APPROVE_URL + '/approve?approvetemplate=' + approveTemplate + '&request=' + row.RequestName)
                .then(response => {
                    if(!response.ok){
                        const error = response.statusText;
                        return Promise.reject(error);
                    }
                    
                    this.setRequestStatus(row)
                })
                .catch(error => {
                    console.error('There was an error!', error);
                });
          } else {
            console.log("Approve failed...")
        }
    }

    setRequestStatus(request) {
        console.log(request)
  
        let approve = this.state.approve
        if(approve) {
          for(let i=0;i<approve.requests.length;i++){
            if(approve.requests[i].RequestName == request.RequestName) {
                approve.requests[i].Approved = true
            }
          }
        }
        
        this.setState({
            approve: approve,
        })
    }

    render(){
      console.log("rendering... hehehe")
      return (
        <Box sx={{ flexGrow: 1}}>
            <AppBar position="static">
                <Toolbar>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                    Tekton Manual Approve
                </Typography>
                </Toolbar>
            </AppBar>
            <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="simple table">
                    <TableHead>
                        <TableRow>
                            <TableCell>Request Name</TableCell>
                            <TableCell align="left">Request Time</TableCell>
                            <TableCell align="left">Is Approved</TableCell>
                            <TableCell align="left">Approved time</TableCell>
                            <TableCell align="left">Approve</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {this.state.approve.requests.map((row) => (
                        <TableRow
                            key={row.RequestName}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.RequestName}
                            </TableCell>
                            <TableCell align="left">{row.RequestTimestamp}</TableCell>
                            <TableCell align="left">{row.Approved?"Yes":"No"}</TableCell>
                            <TableCell align="left">{row.ApproveTimestamp}</TableCell>
                            <TableCell align="left">
                                <IconButton onClick={() => { this.approveRequest(row) }}>
                                    <DoneOutlineIcon color={row.Approved?"disabled":"primary"} />
                                </IconButton>
                            </TableCell>
                        </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>
      )
    }
  }

export default Datatable;
