
using HackSite.Dto;
using JsonConfig;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace HackSite
{
    public static class ParseData
    {
        private static WebCall wc = null;

        public static WebCall APICall
        {
            get
            {
                if (wc == null)
                {
                    wc = new WebCall();
                }
                return wc;
            }
        }

        public static IEnumerable<Dto.UnemploymentData> CallQuandlAPI(DateTime startDate, DateTime endDate)
        {
            List<UnemploymentData> result = new List<UnemploymentData>();
            var startDateString = startDate.ToString("yyyy-MM-dd");
            var endDateString = endDate.ToString("yyyy-MM-dd");
            APICall.APIUri = $"http://localhost:8082/quandl?from={startDateString}&to={endDateString}";

            var rawData = APICall.Get();
            var parsedJson = JsonConfig.Config.ParseJson(rawData);
            var keys = parsedJson.Keys;
            var values = parsedJson.Values;

            var stateData = Newtonsoft.Json.JsonConvert.DeserializeObject<IEnumerable<State>>(State.states);

            foreach (var key in keys)
            {
                var stateName = stateData.Where(x => x.StateAbbr == key).First().StateName;
                foreach (ConfigObject value in (ConfigObject[])parsedJson[key])
                {

                    var data = new UnemploymentData()
                    {
                        State = stateName,
                        StateAbbr = key,
                        Date = ParseDateTime(value.ElementAt(0).Value),
                        Value = ParseDouble(value.ElementAt(1).Value) * 10
                    };

                    result.Add(data);
                }
            }
            return result;

            //var data = JsonConvert.DeserializeObject<List<UnemploymentData>>(rawData);
        }

        public static IEnumerable<NewsData> CallHavenAPI(DateTime startDate)
        {
            List<NewsData> result = new List<NewsData>();
            var startDateString = startDate.ToString("yyyy-MM-dd");

            APICall.APIUri = $"http://localhost:8082/haven?from={startDateString}";
            var rawData = APICall.Get();

            var parsedJson = JsonConfig.Config.ParseJson(rawData);
            var keys = parsedJson.Keys;
            var values = parsedJson.Values;

            var stateData = Newtonsoft.Json.JsonConvert.DeserializeObject<IEnumerable<State>>(State.states);

            foreach (var key in keys)
            {
                var stateName = stateData.Where(x => x.StateAbbr == key).First().StateName;

                foreach (var articleKey in ((ConfigObject)parsedJson[key]).Keys)
                {
                    var valObject = ((ConfigObject)((ConfigObject)parsedJson[key])[articleKey]).Values;

                    var data = new NewsData()
                    {
                        State = stateName,
                        StateAbbr = key,
                        ArticleTitle = articleKey,
                        ArticleLink = valObject.ElementAt(1).ToString(),
                        score = Convert.ToDecimal(((ParseDouble(valObject.ElementAt(0)) * 100).ToString()+"000".Substring(0, 3)))
                    };

                    result.Add(data);

                }
            }

            return result;
        }


        private static DateTime? ParseDateTime(object data)
        {
            DateTime resultData = DateTime.MinValue;
            if (DateTime.TryParse((string)data, out resultData))
            {
                return resultData;
            }
            return null;
        }

        private static decimal? ParseDouble(object data)
        {
            decimal resultData = decimal.MinValue;
            if (decimal.TryParse(data.ToString(), out resultData))
            {
                return resultData;
            }
            return null;
        }
    }
}