{
    'graphite_base': 'http://graphite.thraxil.org/',

    'graphs': [
        {
            'name' : 'base',
            'attrs' : {
                'bgcolor': '#ffffff',
                'colorList': '#999999,#006699',
                'graphOnly': 'true'
            }
        },

        {
            'name': 'servergraph',
            'parents': ['base'],
            'attrs': {
                'width': '300',
                'height': '100'
            }
        },

        {
            'name': 'weekly',
            'attrs': {
                'height': '50',
                'from': '-7days',
                'bgcolor': '#eeeeee',
                'colorList': '#cccccc,#6699cc'
            }
        },

        {
            'name': 'CPU Graph',
            'parents': ['servergraph'],
            'attrs': {
                'yMin': '0',
                'areaMode': 'first'
            },
            'targets': [
                'server.$1.cpu.load_average.1_minute',
                'constantLine(1)'
            ]
        },

        {
            'name': 'Mem Graph',
            'parents': ['servergraph'],
            'attrs': {
                'yMin': '0',
                'yMax': '100',
                'areaMode': 'first'
            },
            'targets': [
                'server.$1.memory.MemFree.percent',
                'constantLine(10)'
            ]
        },

        {
            'name': 'Disk Usage',
            'parents': ['servergraph'],
            'attrs': {
                'yMin': '0',
                'yMax': '100',
                'areaMode': 'first'
            },
            'targets': [
                'asPercent(server.$1.disk.usage.available,server.$1.disk.usage.total)',
                'constantLine(10)'
            ]
        },

        {
            'name': 'Network Graph',
            'parents': ['servergraph'],
            'targets': [
                'nonNegativeDerivative(server.$1.network.eth0.receive.byte_count)',
                'nonNegativeDerivative(server.$1.network.eth0.transmit.byte_count)'
            ]
        },

        {
            'name': 'applicationgraph',
            'parents': ['base'],
            'attrs': {
                'width': '400',
                'height': '100'
            }
        },

        {
            'name': 'reqsgraph',
            'parents': ['applicationgraph'],
            'attrs': {
                'yMin': '0',
                'lineMode': 'connected'
            }
        },

        {
            'name': 'fivehundredsgraph',
            'parents': ['applicationgraph'],
            'attrs': {
                'yMin': '0',
                'drawNullAsZero': 'True',
                'lineMode': 'connected'
            }
        },

        {
            'name': 'timesgraph',
            'parents': ['applicationgraph'],
            'attrs': {
                'yMin': '0',
                'lineMode': 'connected',
            },
            'targets': ['constantLine(100)']
        }

    ],

    'sections': [
        {
            'label': 'Servers',
            'graphs': [
                {
                    'label': ''
                }
            ]
        },
        {
            'label': 'Applications',

        }
    ]
}
