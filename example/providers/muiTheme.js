import { createTheme } from '@mui/material/styles'

const theme = createTheme({
  typography: {
    htmlFontSize: 16,
    fontSize: 14,
    h1: {
      fontSize: 18,
      fontWeight: 'bold',
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          padding: '0.62rem 1.2rem',
          fontWeight: 'bold',
          fontSize: '1rem',
          borderRadius: 24,
          boxShadow: 'none',
          '&:hover': {
            boxShadow: 'none',
          },
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          borderRadius: 6,
        },
      },
    },
    MuiDialogTitle: {
      styleOverrides: {
        root: {
          backgroundColor: '#f1f1f1',
          marginBottom: '1.3rem',
          borderBottom: '1px solid #ddd',
        },
      },
    },
    MuiDialogActions: {
      styleOverrides: {
        root: {
          padding: '.3rem 1rem 1.15rem 0',
          '.MuiButton-root:first-of-type': {
            color: '#555',
          },
        },
      },
    },
    MuiSwitch: {
      styleOverrides: {
        switchBase: {
          color: '#888888',
          '&.Mui-checked': {
            color: '#1066e4',
          },
          '&.Mui-checked + .MuiSwitch-track': {
            backgroundColor: '#7e8af9',
          },
        },
        checked: {},
        track: {
          backgroundColor: '#999',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          minWidth: 380,
          boxShadow: 'none',
          border: '1px solid #ccc',
          backgroundColor: 'transparent',
          '.MuiCardContent-root:last-child': {
            padding: 0,
          },
          '&.danger': {
            border: '1px solid #ff9139',
            backgroundColor: '#fff5f5',
          },
          variants: [
            {
              props: { variant: 'action' },
              style: {
                '> .MuiCardContent-root': {
                  paddingBottom: 0,
                  '.MuiContainer-root': {
                    width: '100%',
                    backgroundColor: 'transparent',
                    '.MuiList-root': {
                      padding: 0,
                      '.MuiListItemIcon-root': {
                        flexDirection: 'row-reverse',
                      },
                      '.MuiListItemText-root:first-of-type .MuiTypography-root': {
                        color: '#666',
                      },
                      '.MuiListItemButton-root .MuiContainer-root': {
                        width: '100%',
                        padding: '14px 0px',
                        paddingTop: 10,
                        '> .MuiTypography-root': {
                          fontSize: '1.2rem',
                          marginBottom: 8,
                        },
                        '.MuiListItemText-root': {
                          margin: 0,
                          '.MuiTypography-root': {
                            fontSize: '1rem',
                            color: 'rgba(0, 0, 0, 0.87)',
                          },
                        },
                      },
                    },
                  },
                },
              },
            },
            {
              props: { variant: 'summary' },
              style: {
                '> .MuiCardContent-root .MuiContainer-root': {
                  '&:first-of-type': {
                    marginTop: 2,
                    padding: 16,
                    '> .MuiTypography-root': {
                      fontSize: '1rem',
                      '&:first-of-type': {
                        fontSize: '1.2rem',
                        marginBottom: 8,
                      },
                    },
                  },
                  '&:last-of-type': {
                    width: '100%',
                    backgroundColor: 'transparent',
                    '.MuiList-root': {
                      paddingBottom: 0,
                      '.MuiListItemButton-root': {
                        padding: 16,
                        '.MuiListItemText-root:first-of-type': {
                          paddingTop: '0.1em',
                          '.MuiTypography-root': {
                            color: '#666',
                            fontSize: '1rem',
                          },
                        },
                      },
                      '.MuiListItemIcon-root': {
                        flexDirection: 'row-reverse',
                      },
                      '.MuiDivider-root:last-of-type': {
                        display: 'none',
                      },
                    },
                  },
                },
                '&.info-card > .MuiCardContent-root .MuiContainer-root:last-of-type .MuiList-root .MuiListItemButton-root .MuiListItemText-root':
                  {
                    minWidth: 140,
                    '&:first-of-type': {
                      maxWidth: 140,
                    },
                    '&:nth-of-type(2)': {
                      whiteSpace: 'nowrap',
                      '> .MuiTypography-root': {
                        textOverflow: 'ellipsis',
                        overflow: 'hidden',
                      },
                    },
                  },
              },
            },
          ],
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          color: '#fff',
          fontSize: '.875rem',
          fontWeight: 'bold',
        },
      },
    },
    MuiTable: {
      styleOverrides: {
        root: {
          variants: [
            {
              props: { variant: 'crosstab' },
              style: {
                width: 'calc(100% - 1px)',
                '> .MuiTableHead-root': {
                  position: 'sticky',
                  top: 0,
                  zIndex: 2,
                  background: '#fff',
                },
                '> .MuiTableBody-root': {
                  border: '1px solid #ddd',
                  '.MuiTableRow-root:last-child': {
                    '.MuiTableCell-root': {
                      borderBottom: 0,
                    },
                  },
                  '.MuiTableRow-root:hover': {
                    backgroundColor: '#f5f5f5',
                  },
                  '.MuiTableCell-root': {
                    borderLeft: '1px solid #ddd',
                    padding: '16px',
                    paddingBottom: '12px',
                  },
                  '.MuiTableCell-root:hover': {
                    backgroundColor: '#f0f0f0',
                  },
                  'th:first-of-type': {
                    borderLeft: 'none',
                    width: '0px',
                    padding: '5px',
                    '> .MuiAvatar-root': {
                      color: '#aaa',
                      backgroundColor: 'transparent',
                      fontSize: '12px',
                    },
                  },
                },
              },
            },
          ],
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        root: {
          flexShrink: 0,
          overflowX: 'hidden',
          width: 240,
          '&.compact': {
            width: 72,
            '> div': {
              width: 72,
              overflowX: 'hidden',
            },
          },
        },
        paper: ({ theme }) => ({
          width: 240,
          backgroundColor: '#f7f7f7',
          '.MuiContainer-root:first-of-type': {
            padding: 0,
            marginBottom: theme.spacing(2),
            '.MuiToolbar-root': {
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'flex-start',
              padding: theme.spacing(0), // TODO (0, 1)
              paddingLeft: theme.spacing(2),
              minHeight: '64px',
              '.MuiIconButton-root': {
                color: 'rgba(0, 0, 0, 0.4)',
              },
            },
          },
          '.MuiContainer-root:last-of-type': {
            padding: 0,
            marginTop: 'auto',
            marginBottom: theme.spacing(1),
          },
          '.MuiList-root': {
            paddingTop: 0,
            a: {
              color: 'rgba(0, 0, 0, 0.6)',
            },
            '.MuiListItemIcon-root': {
              color: 'rgba(0, 0, 0, 0.4)',
            },
            '.Mui-selected': {
              color: 'rgba(0, 0, 0, 0.85)',
            },
            '.Mui-selected .MuiListItemIcon-root': {
              color: 'rgba(0, 0, 0, 0.7)',
            },
            '.MuiListItemButton-root': {
              paddingLeft: theme.spacing(3),
              paddingTop: 12,
              paddingBottom: 12,
              '> .MuiListItemText-root': {
                margin: 0,
              },
            },
            '.MuiListItem-root': {
              padding: 0,
              width: '100%',
              '> .MuiLink-root': {
                width: '100%',
              },
            },
          },
        }),
      },
    },
    MuiBackdrop: {
      styleOverrides: {
        root: ({ theme }) => ({
          variants: [
            {
              props: { variant: 'fullscreen' },
              style: {
                zIndex: theme.zIndex.drawer + 1,
                backgroundColor: 'rgba(255, 255, 255, 0.62)',
                '.MuiLinearProgress-root': {
                  position: 'absolute',
                  top: 0,
                  width: '100%',
                },
              },
            },
          ],
        }),
      },
    },
  },
  palette: {
    background: {
      default: '#fff',
    },
    primary: {
      main: '#000',
    },
    secondary: {
      main: '#0418ef',
    },
  },
})

export default theme
