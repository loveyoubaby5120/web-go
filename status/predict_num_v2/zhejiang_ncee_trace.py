# coding: utf8

from collections import Counter
import random 
import numpy as np
import pandas as pd

def main():
    
    basename = '理一批'
    df_major, df_count = read_param('constrain/' + basename)
    df_major['weight'] = df_major.apply(lambda row: -.5 if row['enroll_num'] in [1,2] else row['avg_score']-row['low_score'], axis=1)
    df_major['sort_high'] = df_major.apply(lambda row: (row['avg_score'] * 1.9 - .2 * row['low_score'])/(row['weight'] + 1), axis=1)
    del df_major['weight']
#     df_major['sort_high'] = df_major.apply(lambda row: (row['avg_score']-row['low_score']) * .9 - .1 * row['enroll_num'], axis=1)
    df_major = df_major.sort_values(by=['low_score','sort_high','enroll_num'],ascending=[False,False,True])
    
#     df_major = pd.read_excel('文一批_r3.xlsx')
    low_scores = Counter(df_major['low_score'])
    scores = Counter(dict(zip(df_count['score'], df_count['count'])))
    
    current = scores - low_scores
    list_solve = []
    for _, (major, enroll_num, avg_score, low_score, rank, sort_high) in df_major.iterrows():
        if enroll_num > 1:
            lower_sum = (avg_score - .48) * enroll_num - low_score
            upper_sum = (avg_score + .48) * enroll_num - low_score
            
            pool = create_pool(current, low_score)
            flag, solution = find_solution(enroll_num-1, lower_sum, upper_sum, pool, [])
            print(solution)
            
            current = current - Counter(solution)
            solution.append(low_score)
            list_solve.append([major, enroll_num, avg_score, low_score, rank, sort_high, flag, solution])
        else:
            list_solve.append([major, enroll_num, avg_score, low_score, rank, sort_high, 'only 1', [low_score]])
        
    df_solve = pd.DataFrame(list_solve,
                            columns=['major_key','enroll_num','avg_score','low_score','rank','sort_high','flag','combination'])
    df_solve.to_excel('{}_solution.xlsx'.format(basename), index=False)
    
    save_counter(scores, '{}_origin_scores.xlsx'.format(basename))
    save_counter(current, '{}_rest_scores.xlsx'.format(basename))
    print('ok')

def save_counter(counter, path):
    
    df_data = pd.DataFrame.from_dict(counter, orient='index')
    df_data.to_excel(path)

def read_param(basename):
    
    df_major = pd.read_excel('{}.xlsx'.format(basename), sheetname='高考分数')
    df_major['major'] = df_major.apply(lambda row: '{}@{}'.format(row['学校'],row['专业']), axis=1)
    df_major.rename(columns=dict(zip(['录取人数','平均分','最低分','名次号'],
                                     ['enroll_num','avg_score','low_score','rank'])), inplace=True)
    
    df_count = pd.read_excel('{}.xlsx'.format(basename), sheetname='一分一段')
    df_count['count'] = df_count['小计/累计'].apply(lambda x: int(x.split('/')[0]))
    df_count.rename(columns={'分数':'score'}, inplace=True)
    
    return (df_major.loc[:,['major','enroll_num','avg_score','low_score','rank']], 
            df_count.loc[:, ['score','count']])
    
def create_pool(counter, lbound):
    list_ = sorted(counter.elements(), reverse=True)
    return np.array(list(filter(lambda x: x >= lbound, list_)))
    

def find_solution(n, lower_sum, upper_sum, pool, result=[]):
    
    if n == 0: 
        return 'blank', result
    
    if len(pool) < n:
        print('Not enough students left to choose!')
        return  'not enough', result
    
    if len(pool) == n or len(set(pool)) == 1:
        if lower_sum <= sum(pool[:n]) <= upper_sum:
            return 'ok', result + list(pool[:n])
        else:
            print('Just n scores left or ONLY one unique value: but the sum is out of the boundary!')
            return  'even value', result
        
    if len(pool) > n:
        # max/min N-numbers out of the boundary
        if sum(pool[:n]) < lower_sum:
            print('Sum of max N-numbers is less than the lower bound!') 
            return 'pool too small', result
        if sum(pool[-n:]) > upper_sum:
            print('Sum of min N-numbers is greater than the upper bound!') 
            return 'pool too large', result
        
        # general procedure
        if n == 1:
            upper_b = upper_sum
        else:
            print(n, pool[:10])
            upper_b = upper_sum - sum(pool[-n+1:])
        
        i = sum( (pool - upper_b) > 0)  # how many numbers are greater than the MAX-maybe
        guess = pool[i]
        result.append(guess)
        
        if n == 1: return 'ok', result
        return find_solution(n-1, lower_sum-guess, upper_sum-guess, pool[i+1:], result)
            

if __name__ == '__main__': main()