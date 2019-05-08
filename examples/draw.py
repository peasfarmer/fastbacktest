# coding=utf-8

# coding=utf8
import json
import msgpack
import dash
import dash_core_components as dcc
import dash_html_components as html
from dash.dependencies import Input, Output
import pandas as pd
import plotly.graph_objs as go
import plotly
import threading
import os, time, sys
import gc

g_dfList = []

MINTS = pd.to_datetime('1970-1-1')
MAXTS = pd.to_datetime('2200-1-1')

app = dash.Dash(__name__)

app.layout = html.Div([
    dcc.Graph(
        id='basic-interactions',
        style={"height": "100vh"}
    ),

])


@app.callback(
    dash.dependencies.Output('basic-interactions', 'figure'),
    [Input('basic-interactions', 'relayoutData')],
    state=[dash.dependencies.State('basic-interactions', 'figure')]
)
def display_selected_data(relayoutData, state):
    hiddenList = []
    if state is not None:
        lines = state.get('data', {})
        for line in lines:
            if line.get('visible', 'legendonly') == 'legendonly':
                hiddenList.append(line['name'])

    if relayoutData is None:
        return showDF(MINTS, MAXTS, hiddenList)
    elif relayoutData.get('autosize', False):
        return showDF(MINTS, MAXTS, hiddenList)
    else:
        start = relayoutData.get('xaxis.range[0]', MINTS)
        if type(start) == str:
            start = pd.to_datetime(start)
        end = relayoutData.get('xaxis.range[1]', MAXTS)
        if type(end) == str:
            end = pd.to_datetime(end)
        return showDF(start, end, hiddenList)


def showDF(start, end, hiddenList):
    global g_dfList

    rows = max([info['plotIndex'] for info in g_dfList])
    fig = plotly.tools.make_subplots(rows=rows, cols=1, specs=[[{}]] * rows,
                                     shared_xaxes=True, shared_yaxes=True, row_width=[i for i in range(1,rows+1)],
                                     vertical_spacing=0.001)

    for dfInfo in g_dfList:
        index = dfInfo.get('plotIndex', 1)
        mode = dfInfo.get('mode', 'lines')
        marker = dfInfo.get('marker', {})
        line = dfInfo.get('line', {})
        text = dfInfo.get('text', None)

        tmp = dfInfo['df']
        if len(tmp) == 0:
            continue
        tmp = limitDFNum(tmp[(tmp.index > start) & (tmp.index < end)])
        for colName in dfInfo['df'].columns:
            visible = True
            if colName in hiddenList:
                visible = 'legendonly'
            trace = go.Scatter(
                x=tmp.index,
                name=colName,
                y=tmp[colName],
                visible=visible,
                mode=mode,
                text=text,
                marker=marker,
                line=line
            )
            fig.append_trace(trace, index, 1)

    # height=800, width=1080,
    fig['layout'].update(xaxis=dict(
        tickformat="%y-%m-%d %H:%M:%S",
        # rangeselector=dict(),
        # rangeslider=dict(), #选择区域活动条
        # type='date'
    ))
    return fig


# 样本量太大的时候抽样显示就可以了
def limitDFNum(df):
    if len(df) < 3000:  # 防止是0长度
        return df

    cols = {}
    isNeedResample = False
    for c in df.columns:
        if c.find("diff") == -1:
            cols[c] = "avg"
        else:
            if c.find('buy') >= 0:
                cols[c] = "min"
                isNeedResample = True
            elif c.find('sell') >= 0:
                cols[c] = "max"
                isNeedResample = True
            else:
                cols[c] = "avg"

    if isNeedResample:
        timeLong = (df.index[-1] - df.index[0]).total_seconds()
        rule = '10T'
        if timeLong > 15 * 24 * 3600:
            rule = '20T'
        elif timeLong > 5 * 24 * 3600:
            rule = '5T'
        else:
            rule = '1T'

        tmp = df.resample(rule).agg(cols).sort_index()
        # print tmp
        return tmp
    else:
        frac = 3000.0 / len(df)  # 抽样百分比
        if frac > 0.9:
            return df
        return df.sample(frac=frac).copy().sort_index()


# def draw(firstDFList, nameList, secondDFList):
#     global g_firstDFList, g_NameList, g_secondDFList
#     g_firstDFList = firstDFList
#     g_NameList = nameList
#     g_secondDFList = secondDFList
#
#     app.run_server(debug=False)
#     pass

def draw():
    # global g_dfList
    # g_dfList = dfList

    app.run_server(debug=False)
    pass


def getDrawFileList():
    ret = []
    path = ''
    if len(sys.argv) == 2:
        path = sys.argv[1]
    else:
        path = path + './draw_data/'
    for f in os.listdir(path):
        subpath = path + '/' + f
        if not os.path.isfile(subpath):
            ret.append(subpath)

    return ret


def msgp2df(fileName):
    data = open(fileName, 'rb').read()
    data2 = msgpack.unpackb(data, raw=False)
    df = pd.DataFrame(data=data2)
    df.index = pd.to_datetime(df['Ts'], unit='ms')
    del df['Ts']

    return df


def loadOneData(path):
    try:
        config = open(path + '/config.json').read()
        config = json.loads(config)
        if config is None:
            return None
        if not config.get('enable'):
            return None
        gc.disable()
        df = msgp2df(path + '/' + config['fileName'])
        gc.enable()
        config.pop('fileName')
        if 'text' in df.columns:
            config.update({
                'text': df['text'].values
            })
            del df['text']
        config.update({
            'df': df
        })
        return config
    except:
        return None


def getDirLastModifyTime():
    lastTS = 0
    drawFiles = getDrawFileList()
    for filePath in drawFiles:
        filePath = filePath + '/config.json'
        t = os.path.getmtime(filePath)
        if t > lastTS:
            lastTS = t
    return lastTS


if __name__ == '__main__':
    lastTS = 0
    lastUpdate = 0

    t1 = threading.Thread(target=draw, args=())
    t1.start()
    while True:
        ts = getDirLastModifyTime()
        # print(datetime.datetime.fromtimestamp(ts).strftime("%Y-%m-%d %H:%M:%S"))
        if ts == lastTS:
            time.sleep(1)
            continue
        else:
            if time.time() - ts < 3:  # 防止数据没有写完就去加载
                time.sleep(1)
                continue
        print("loading...")
        lastTS = ts
        tmp = []
        time.sleep(1)  # 让数据写完
        drawFiles = getDrawFileList()
        for filePath in drawFiles:
            info = loadOneData(filePath)
            if info is not None:
                tmp.append(info)
        g_dfList = tmp
        print(time.strftime("%Y-%m-%d %H:%M:%S"))
